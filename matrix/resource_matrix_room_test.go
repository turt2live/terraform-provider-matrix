package matrix

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"fmt"
	"github.com/turt2live/terraform-provider-matrix/matrix/api"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"regexp"
	"net/http"
	"strings"
	"net/url"
)

type testAccMatrixRoom struct {
	RoomId        string
	CreatorUserId string
	Name          string
	AvatarMxc     string
	Topic         string
	GuestsAllowed bool
	CreatorToken  string
	Preset        string
}

func testAccCreateMatrixRoom(name string, avatarMxc string, topic string, guestsAllowed bool, preset string) (*testAccMatrixRoom) {
	guestAccess := api.RoomGuestAccessEventContent{Policy: "forbidden"}
	if guestsAllowed {
		guestAccess.Policy = "can_join"
	}

	request := &api.CreateRoomRequest{
		Preset: preset,
		InitialState: []api.CreateRoomStateEvent{
			{Type: "m.room.name", Content: api.RoomNameEventContent{Name: name}},
			{Type: "m.room.avatar", Content: api.RoomAvatarEventContent{AvatarMxc: avatarMxc}},
			{Type: "m.room.topic", Content: api.RoomTopicEventContent{Topic: topic}},
			{Type: "m.room.guest_access", Content: guestAccess},
		},
	}

	response := &api.RoomIdResponse{}
	urlStr := api.MakeUrl(testAccClientServerUrl(), "/_matrix/client/r0/createRoom")
	err := api.DoRequest("POST", urlStr, request, response, testAccAdminToken())
	if err != nil {
		panic(err)
	}

	creatorResponse := &api.RoomCreateEventContent{}
	urlStr = api.MakeUrl(testAccClientServerUrl(), "/_matrix/client/r0/rooms", response.RoomId, "/state/m.room.create")
	err = api.DoRequest("GET", urlStr, nil, creatorResponse, testAccAdminToken())
	if err != nil {
		panic(err)
	}

	return &testAccMatrixRoom{
		RoomId:        response.RoomId,
		CreatorUserId: creatorResponse.CreatorUserId,
		Name:          name,
		AvatarMxc:     avatarMxc,
		Topic:         topic,
		GuestsAllowed: guestsAllowed,
		CreatorToken:  testAccAdminToken(),
		Preset:        preset,
	}
}

var testAccMatrixRoomConfig_existingRoom = `
resource "matrix_room" "foobar" {
	room_id = "%s"
	member_access_token = "%s"
}`

func TestAccMatrixRoom_ExistingRoom(t *testing.T) {
	room := testAccCreateMatrixRoom("Sample", "mxc://localhost/AvatarHere", "This is a topic", true, "private_chat")
	conf := fmt.Sprintf(testAccMatrixRoomConfig_existingRoom, room.RoomId, room.CreatorToken)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMatrixRoomDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixRoomExists("matrix_room.foobar"),
					testAccCheckMatrixRoomIdMatchesRoomId("matrix_room.foobar"),
					testAccCheckMatrixRoomMatchesCreated("matrix_room.foobar", room),
					resource.TestCheckResourceAttr("matrix_room.foobar", "id", room.RoomId),
					resource.TestCheckResourceAttr("matrix_room.foobar", "creator_user_id", room.CreatorUserId),
					resource.TestCheckResourceAttr("matrix_room.foobar", "member_access_token", room.CreatorToken),
					resource.TestCheckResourceAttr("matrix_room.foobar", "room_id", room.RoomId),
					resource.TestCheckNoResourceAttr("matrix_room.foobar", "preset"),
					resource.TestCheckResourceAttr("matrix_room.foobar", "name", room.Name),
					resource.TestCheckResourceAttr("matrix_room.foobar", "avatar_mxc", room.AvatarMxc),
					resource.TestCheckResourceAttr("matrix_room.foobar", "topic", room.Topic),
					resource.TestCheckNoResourceAttr("matrix_room.foobar", "invite_user_ids"),
					resource.TestCheckNoResourceAttr("matrix_room.foobar", "local_alias_localpart"),
					resource.TestCheckResourceAttr("matrix_room.foobar", "guests_allowed", strconv.FormatBool(room.GuestsAllowed)),
				),
			},
		},
	})
}

var testAccMatrixRoomConfig_newRoom = `
resource "matrix_room" "foobar" {
	creator_user_id = "%s"
	member_access_token = "%s"
	name = "%s"
	avatar_mxc = "%s"
	topic = "%s"
	guests_allowed = %t
	preset = "%s"
}`

func TestAccMatrixRoom_NewRoom(t *testing.T) {
	creator := testAccCreateTestUser("test_acc_room_create_new")
	room := &testAccMatrixRoom{
		CreatorToken:  creator.AccessToken,
		CreatorUserId: creator.UserId,
		GuestsAllowed: true,
		Topic:         "This is a topic",
		AvatarMxc:     "mxc://localhost/AvatarHere",
		Name:          "Sample",
		Preset:        "private_chat",
	}
	conf := fmt.Sprintf(testAccMatrixRoomConfig_newRoom, room.CreatorUserId, room.CreatorToken, room.Name, room.AvatarMxc, room.Topic, room.GuestsAllowed, room.Preset)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMatrixRoomDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixRoomExists("matrix_room.foobar"),
					testAccCheckMatrixRoomIdMatchesRoomId("matrix_room.foobar"),
					resource.TestMatchResourceAttr("matrix_room.foobar", "id", regexp.MustCompile("^!.+")),
					resource.TestCheckResourceAttr("matrix_room.foobar", "creator_user_id", room.CreatorUserId),
					resource.TestCheckResourceAttr("matrix_room.foobar", "member_access_token", room.CreatorToken),
					resource.TestCheckResourceAttr("matrix_room.foobar", "preset", room.Preset),
					resource.TestCheckResourceAttr("matrix_room.foobar", "name", room.Name),
					resource.TestCheckResourceAttr("matrix_room.foobar", "avatar_mxc", room.AvatarMxc),
					resource.TestCheckResourceAttr("matrix_room.foobar", "topic", room.Topic),
					resource.TestCheckNoResourceAttr("matrix_room.foobar", "invite_user_ids"),
					resource.TestCheckNoResourceAttr("matrix_room.foobar", "local_alias_localpart"),
					resource.TestCheckResourceAttr("matrix_room.foobar", "guests_allowed", strconv.FormatBool(room.GuestsAllowed)),
				),
			},
		},
	})
}

var testAccMatrixRoomConfig_invites = `
resource "matrix_room" "foobar" {
	creator_user_id = "%s"
	member_access_token = "%s"
	invite_user_ids = ["%s", "%s"]
}`

func TestAccMatrixRoom_Invites(t *testing.T) {
	creator := testAccCreateTestUser("test_acc_room_invites")
	targetA := testAccCreateTestUser("test_acc_room_invites_user_a")
	targetB := testAccCreateTestUser("test_acc_room_invites_user_b")
	inviteUserIds := []string{targetA.UserId, targetB.UserId}
	room := &testAccMatrixRoom{
		CreatorToken:  creator.AccessToken,
		CreatorUserId: creator.UserId,
		GuestsAllowed: true,
		Topic:         "This is a topic",
		AvatarMxc:     "mxc://localhost/AvatarHere",
		Name:          "Sample",
		Preset:        "private_chat",
	}
	conf := fmt.Sprintf(testAccMatrixRoomConfig_invites, room.CreatorUserId, room.CreatorToken, inviteUserIds[0], inviteUserIds[1])

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMatrixRoomDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixRoomExists("matrix_room.foobar"),
					testAccCheckMatrixRoomIdMatchesRoomId("matrix_room.foobar"),
					testAccCheckMatrixRoomInvitedUsers("matrix_room.foobar", inviteUserIds),
					resource.TestMatchResourceAttr("matrix_room.foobar", "id", regexp.MustCompile("^!.+")),
					resource.TestCheckResourceAttr("matrix_room.foobar", "creator_user_id", room.CreatorUserId),
					resource.TestCheckResourceAttr("matrix_room.foobar", "member_access_token", room.CreatorToken),
				),
			},
		},
	})
}

var testAccMatrixRoomConfig_localAlias = `
resource "matrix_room" "foobar" {
	creator_user_id = "%s"
	member_access_token = "%s"
	local_alias_localpart = "%s"
}`

func TestAccMatrixRoom_LocalAlias(t *testing.T) {
	creator := testAccCreateTestUser("test_acc_room_local_alias")
	room := &testAccMatrixRoom{
		CreatorToken:  creator.AccessToken,
		CreatorUserId: creator.UserId,
		GuestsAllowed: true,
		Topic:         "This is a topic",
		AvatarMxc:     "mxc://localhost/AvatarHere",
		Name:          "Sample",
		Preset:        "private_chat",
	}
	expectedAlias := "test_acc_room_local_alias"
	conf := fmt.Sprintf(testAccMatrixRoomConfig_localAlias, room.CreatorUserId, room.CreatorToken, expectedAlias)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMatrixRoomDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixRoomExists("matrix_room.foobar"),
					testAccCheckMatrixRoomIdMatchesRoomId("matrix_room.foobar"),
					testAccCheckMatrixRoomLocalAliasMatches("matrix_room.foobar", expectedAlias, room.CreatorUserId),
					resource.TestMatchResourceAttr("matrix_room.foobar", "id", regexp.MustCompile("^!.+")),
					resource.TestCheckResourceAttr("matrix_room.foobar", "creator_user_id", room.CreatorUserId),
					resource.TestCheckResourceAttr("matrix_room.foobar", "member_access_token", room.CreatorToken),
					resource.TestCheckResourceAttr("matrix_room.foobar", "local_alias_localpart", expectedAlias),
				),
			},
		},
	})
}

func testAccCheckMatrixRoomDestroy(s *terraform.State) error {
	// TODO: Check that the room was ""deleted""
	return nil
}

func testAccCheckMatrixRoomExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		memberToken := rs.Primary.Attributes["member_access_token"]

		// We'll try to query something like the create event to prove the room exists
		response := &api.RoomCreateEventContent{}
		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms/", rs.Primary.ID, "/state/m.room.create")
		err := api.DoRequest("GET", urlStr, nil, response, memberToken)
		if err != nil {
			return err
		}

		if response.CreatorUserId == "" {
			return fmt.Errorf("creator user id is empty")
		}

		return nil
	}
}

func testAccCheckMatrixRoomMatchesCreated(n string, room *testAccMatrixRoom) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		memberAccessToken := rs.Primary.Attributes["member_access_token"]
		roomId := rs.Primary.ID

		nameResponse := &api.RoomNameEventContent{}
		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms/", roomId, "/state/m.room.name")
		err := api.DoRequest("GET", urlStr, nil, nameResponse, memberAccessToken)
		if err != nil {
			if r, ok := err.(*api.ErrorResponse); !ok || r.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error getting room name: %s", err)
			}
		}

		avatarResponse := &api.RoomAvatarEventContent{}
		urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms/", roomId, "/state/m.room.avatar")
		err = api.DoRequest("GET", urlStr, nil, avatarResponse, memberAccessToken)
		if err != nil {
			if r, ok := err.(*api.ErrorResponse); !ok || r.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error getting room avatar: %s", err)
			}
		}

		topicResponse := &api.RoomTopicEventContent{}
		urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms/", roomId, "/state/m.room.topic")
		err = api.DoRequest("GET", urlStr, nil, topicResponse, memberAccessToken)
		if err != nil {
			if r, ok := err.(*api.ErrorResponse); !ok || r.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error getting room topic: %s", err)
			}
		}

		guestResponse := &api.RoomGuestAccessEventContent{}
		urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms/", roomId, "/state/m.room.guest_access")
		err = api.DoRequest("GET", urlStr, nil, guestResponse, memberAccessToken)
		if err != nil {
			if r, ok := err.(*api.ErrorResponse); !ok || r.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error getting room guest access policy: %s", err)
			}
		}

		creatorResponse := &api.RoomCreateEventContent{}
		urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms/", roomId, "/state/m.room.create")
		err = api.DoRequest("GET", urlStr, nil, creatorResponse, memberAccessToken)
		if err != nil {
			if r, ok := err.(*api.ErrorResponse); !ok || r.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error getting room creator: %s", err)
			}
		}

		joinRulesResponse := &api.RoomJoinRulesEventContent{}
		urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms/", roomId, "/state/m.room.join_rules")
		err = api.DoRequest("GET", urlStr, nil, joinRulesResponse, memberAccessToken)
		if err != nil {
			if r, ok := err.(*api.ErrorResponse); !ok || r.StatusCode != http.StatusNotFound {
				return fmt.Errorf("error getting room join rule policy: %s", err)
			}
		}

		if creatorResponse.CreatorUserId != room.CreatorUserId {
			return fmt.Errorf("creator mismatch. expected: %s  got: %s", room.CreatorUserId, creatorResponse.CreatorUserId)
		}
		if nameResponse.Name != room.Name {
			return fmt.Errorf("name mismatch. expected: %s  got: %s", room.Name, nameResponse.Name)
		}
		if avatarResponse.AvatarMxc != room.AvatarMxc {
			return fmt.Errorf("avatar mismatch. expected: %s  got: %s", room.AvatarMxc, avatarResponse.AvatarMxc)
		}
		if topicResponse.Topic != room.Topic {
			return fmt.Errorf("topic mismatch. expected: %s  got: %s", room.Topic, topicResponse.Topic)
		}

		guestPolicy := "forbidden"
		if room.GuestsAllowed {
			guestPolicy = "can_join"
		}
		if guestResponse.Policy != guestPolicy {
			return fmt.Errorf("guest_access mismatch. expected: %s  got: %s", guestPolicy, guestResponse.Policy)
		}

		joinPolicy := "public"
		if room.Preset == "private_chat" {
			joinPolicy = "invite"
		}
		if joinRulesResponse.Policy != joinPolicy {
			return fmt.Errorf("join_rules mismatch. expected: %s  got: %s", joinPolicy, joinRulesResponse.Policy)
		}

		return nil
	}
}

func testAccCheckMatrixRoomIdMatchesRoomId(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		if rs.Primary.ID != rs.Primary.Attributes["room_id"] {
			return fmt.Errorf("room_id and record id do not match")
		}

		return nil
	}
}

func testAccCheckMatrixRoomInvitedUsers(n string, invitedUserIds []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		memberAccessToken := rs.Primary.Attributes["member_access_token"]
		roomId := rs.Primary.ID

		for _, invitedUserId := range invitedUserIds {
			response := &api.RoomMemberEventContent{}
			urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/rooms/", roomId, "/state/m.room.member/", invitedUserId)
			err := api.DoRequest("GET", urlStr, nil, response, memberAccessToken)
			if err != nil {
				return fmt.Errorf("error getting room member %s: %s", invitedUserId, err)
			}
			if response.Membership != "invite" {
				return fmt.Errorf("user %s is not invited", invitedUserId)
			}
		}

		return nil
	}
}

func testAccCheckMatrixRoomLocalAliasMatches(n string, aliasLocalpart string, creatorUserId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		memberAccessToken := rs.Primary.Attributes["member_access_token"]
		roomId := rs.Primary.ID

		// We're forced to do an estimation on what the full alias will look like, so we try and get the
		// homeserver domain from the creator's user ID. This is bad practice, however most of the tests
		// will be run on localhost anyways, making this check not so bad.
		idParts := strings.Split(creatorUserId, ":")
		if len(idParts) != 2 && len(idParts) != 3 {
			return fmt.Errorf("illegal matrix user id: %s", creatorUserId)
		}
		hsDomain := idParts[1]
		if len(idParts) > 2 { // port
			hsDomain = fmt.Sprintf("%s:%s", hsDomain, idParts[2])
		}
		fullAlias := fmt.Sprintf("#%s:%s", aliasLocalpart, hsDomain)
		safeAlias := url.QueryEscape(fullAlias)

		response := &api.RoomDirectoryLookupResponse{}
		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/directory/room/", safeAlias)
		err := api.DoRequest("GET", urlStr, nil, response, memberAccessToken)
		if err != nil {
			return fmt.Errorf("error querying alias: %s", err)
		}
		if response.RoomId != roomId {
			return fmt.Errorf("room id mismatch. expected: %s  got: %s", roomId, response.RoomId)
		}

		return nil
	}
}
