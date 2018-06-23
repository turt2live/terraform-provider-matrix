package api

type RoomMemberEvent struct {
	Content  *RoomMemberEventContent `json:"content"`
	Type     string                  `json:"type"`
	EventId  string                  `json:"event_id"`
	RoomId   string                  `json:"room_id"`
	StateKey string                  `json:"state_key"`

	// other fields not included
}
