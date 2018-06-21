package matrix

import (
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
	"log"
	"github.com/turt2live/terraform-provider-matrix/matrix/api"
	"net/http"
)

func resourceRoom() *schema.Resource {
	return &schema.Resource{
		Exists: resourceRoomExists,
		Create: resourceRoomCreate,
		Read:   resourceRoomRead,
		Update: resourceRoomUpdate,
		Delete: resourceRoomDelete,

		Schema: map[string]*schema.Schema{
			"creator_user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"member_access_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"room_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"preset": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// Ignored if no creator
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"avatar_mxc": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"invite_user_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
				// Ignored if no creator
			},
			"local_alias_localpart": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// Ignored if no creator
			},
			"guests_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceRoomCreate(d *schema.ResourceData, m interface{}) error {
	meta := m.(Metadata)

	creatorIdRaw := nilIfEmptyString(d.Get("creator_user_id"))
	memberAccessToken := d.Get("member_access_token").(string)
	roomIdRaw := nilIfEmptyString(d.Get("room_id"))

	presetRaw := d.Get("preset").(string)
	nameRaw := nilIfEmptyString(d.Get("name"))
	avatarMxcRaw := nilIfEmptyString(d.Get("avatar_mxc"))
	topicRaw := nilIfEmptyString(d.Get("topic"))
	aliasLocalpartRaw := d.Get("local_alias_localpart").(string)
	guestsAllowed := d.Get("guests_allowed").(bool)
	invitedUserIds := d.Get("invite_user_ids").([]string)

	hasCreator := creatorIdRaw != nil
	hasRoomId := roomIdRaw != nil

	if hasCreator && hasRoomId {
		return fmt.Errorf("cannot specify both a creator and room_id")
	}

	if !hasCreator && !hasRoomId {
		return fmt.Errorf("a creator or room_id must be specified")
	}

	if hasCreator {
		log.Printf("[DEBUG] Creating room")
		request := &api.CreateRoomRequest{
			Preset:         presetRaw,
			AliasLocalpart: aliasLocalpartRaw,
			AllowGuests:    guestsAllowed,
			InviteUserIds:  invitedUserIds,
		}

		stateEvents := make([]api.CreateRoomStateEvent, 0)
		if nameRaw != nil {
			stateEvents = append(stateEvents, api.CreateRoomStateEvent{
				Type:    "m.room.name",
				Content: api.RoomNameEventContent{Name: nameRaw.(string)},
			})
		}
		if avatarMxcRaw != nil {
			stateEvents = append(stateEvents, api.CreateRoomStateEvent{
				Type:    "m.room.avatar",
				Content: api.RoomAvatarEventContent{AvatarMxc: avatarMxcRaw.(string)},
			})
		}
		if topicRaw != nil {
			stateEvents = append(stateEvents, api.CreateRoomStateEvent{
				Type:    "m.room.topic",
				Content: api.RoomTopicEventContent{Topic: topicRaw.(string)},
			})
		}
		request.InitialState = stateEvents

		response := &api.RoomIdResponse{}
		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/createRoom")
		err := api.DoRequest("POST", urlStr, request, response, memberAccessToken)
		if err != nil {
			return fmt.Errorf("error creating room: %s", err)
		}

		d.SetId(response.RoomId)
		d.Set("room_id", response.RoomId)
	}

	return resourceRoomRead(d, meta)
}

func resourceRoomExists(d *schema.ResourceData, m interface{}) (bool, error) {
	meta := m.(Metadata)

	memberAccessToken := d.Get("member_access_token").(string)
	roomIdRaw := nilIfEmptyString(d.Get("room_id"))

	if roomIdRaw == nil {
		return false, nil
	}

	// First identify who the user is
	log.Printf("[DEBUG] Doing whoami on: %s ", d.Id())
	urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/account/whoami")
	whoAmIResponse := &api.WhoAmIResponse{}
	err := api.DoRequest("GET", urlStr, nil, whoAmIResponse, memberAccessToken)
	if err != nil {
		// We say true so that Terraform won't accidentally delete the room
		return true, fmt.Errorf("error performing whoami: %s", err)
	}

	// Now that we have user's ID, let's make sure they are a member
	memberEventResponse := &api.RoomMemberEventContent{}
	urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms", roomIdRaw.(string), "/state/m.room.member/", whoAmIResponse.UserId)
	err = api.DoRequest("GET", urlStr, nil, memberEventResponse, memberAccessToken)
	if err != nil {
		// An error accessing the room means it doesn't exist anymore
		return false, fmt.Errorf("error getting member event for user: %s", err)
	}

	if memberEventResponse.Membership != "join" {
		return false, fmt.Errorf("member is not in the room")
	}

	return true, nil
}

func resourceRoomRead(d *schema.ResourceData, m interface{}) error {
	meta := m.(Metadata)

	memberAccessToken := d.Get("member_access_token").(string)
	roomIdRaw := nilIfEmptyString(d.Get("room_id"))

	if roomIdRaw == nil {
		return fmt.Errorf("no room_id")
	}

	nameResponse := &api.RoomNameEventContent{}
	urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms", roomIdRaw.(string), "/state/m.room.name")
	err := api.DoRequest("GET", urlStr, nil, nameResponse, memberAccessToken)
	if err != nil {
		if r, ok := err.(*api.ErrorResponse); !ok || r.StatusCode != http.StatusNotFound {
			return fmt.Errorf("error getting room name: %s", err)
		}
	}

	avatarResponse := &api.RoomAvatarEventContent{}
	urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms", roomIdRaw.(string), "/state/m.room.avatar")
	err = api.DoRequest("GET", urlStr, nil, avatarResponse, memberAccessToken)
	if err != nil {
		if r, ok := err.(*api.ErrorResponse); !ok || r.StatusCode != http.StatusNotFound {
			return fmt.Errorf("error getting room avatar: %s", err)
		}
	}

	topicResponse := &api.RoomTopicEventContent{}
	urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms", roomIdRaw.(string), "/state/m.room.topic")
	err = api.DoRequest("GET", urlStr, nil, topicResponse, memberAccessToken)
	if err != nil {
		if r, ok := err.(*api.ErrorResponse); !ok || r.StatusCode != http.StatusNotFound {
			return fmt.Errorf("error getting room topic: %s", err)
		}
	}

	guestResponse := &api.RoomGuestAccessEventContent{}
	urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms", roomIdRaw.(string), "/state/m.room.guest_access")
	err = api.DoRequest("GET", urlStr, nil, guestResponse, memberAccessToken)
	if err != nil {
		if r, ok := err.(*api.ErrorResponse); !ok || r.StatusCode != http.StatusNotFound {
			return fmt.Errorf("error getting room guest access policy: %s", err)
		}
	}

	d.Set("name", nameResponse.Name)
	d.Set("avatar_mxc", avatarResponse.AvatarMxc)
	d.Set("topic", topicResponse.Topic)

	if guestResponse.Policy == "can_join" {
		d.Set("guests_allowed", true)
	} else {
		d.Set("guests_allowed", false)
	}

	return nil
}

func resourceRoomUpdate(d *schema.ResourceData, m interface{}) error {
	meta := m.(Metadata)

	memberAccessToken := d.Get("member_access_token").(string)
	roomIdRaw := nilIfEmptyString(d.Get("room_id"))

	if roomIdRaw == nil {
		return fmt.Errorf("no room_id")
	}

	if d.HasChange("name") {
		request := &api.RoomNameEventContent{Name: d.Get("name").(string)}
		response := &api.EventIdResponse{}
		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms", roomIdRaw.(string), "/state/m.room.name")
		err := api.DoRequest("PUT", urlStr, request, response, memberAccessToken)
		if err != nil {
			return err
		}
	}

	if d.HasChange("avatar_mxc") {
		request := &api.RoomAvatarEventContent{AvatarMxc: d.Get("avatar_mxc").(string)}
		response := &api.EventIdResponse{}
		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms", roomIdRaw.(string), "/state/m.room.avatar")
		err := api.DoRequest("PUT", urlStr, request, response, memberAccessToken)
		if err != nil {
			return err
		}
	}

	if d.HasChange("topic") {
		request := &api.RoomTopicEventContent{Topic: d.Get("topic").(string)}
		response := &api.EventIdResponse{}
		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms", roomIdRaw.(string), "/state/m.room.topic")
		err := api.DoRequest("PUT", urlStr, request, response, memberAccessToken)
		if err != nil {
			return err
		}
	}

	if d.HasChange("guests_allowed") {
		policy := "forbidden"
		if d.Get("guests_allowed").(bool) {
			policy = "can_join"
		}
		request := &api.RoomGuestAccessEventContent{Policy: policy}
		response := &api.EventIdResponse{}
		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms", roomIdRaw.(string), "/state/m.room.guest_access")
		err := api.DoRequest("PUT", urlStr, request, response, memberAccessToken)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceRoomDelete(d *schema.ResourceData, m interface{}) error {
	// Rooms cannot technically be deleted, so we just abandon them instead
	// Abandoning means kicking everyone and leaving it to rot away. Before we leave though, we'll make sure no one can
	// get back in.

	// TODO: Remove aliases
	// TODO: Set invite only
	// TODO: Disable guest access
	// TODO: Close history to members only
	// TODO: Kick everyone
	// TODO: Leave

	return nil
}
