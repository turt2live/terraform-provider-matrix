package api

type RoomNameEventContent struct {
	Name string `json:"name"`
}

type RoomTopicEventContent struct {
	Topic string `json:"topic"`
}

type RoomAvatarEventContent struct {
	AvatarMxc string `json:"url"`
}

type RoomMemberEventContent struct {
	DisplayName string `json:"displayname,omitempty"`
	AvatarMxc   string `json:"avatar_url,omitempty"`
	Membership  string `json:"membership"`
}

type RoomGuestAccessEventContent struct {
	Policy string `json:"guest_access"`
}
