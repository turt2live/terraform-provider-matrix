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

type RoomCreateEventContent struct {
	CreatorUserId string `json:"creator"`
}

type RoomJoinRulesEventContent struct {
	Policy string `json:"join_rule"`
}

type RoomAliasesEventContent struct {
	Aliases []string `json:"aliases,flow"`
}
