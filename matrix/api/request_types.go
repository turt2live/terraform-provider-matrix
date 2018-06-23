package api

type RegisterRequest struct {
	Authentication           *RegisterAuthenticationData `json:"auth,omitempty"`
	BindEmail                bool                        `json:"bind_email,omitempty"`
	Username                 string                      `json:"username,omitempty"`
	Password                 string                      `json:"password,omitempty"`
	DeviceId                 string                      `json:"device_id,omitempty"`
	InitialDeviceDisplayName string                      `json:"initial_device_display_name,omitempty"`
}

type RegisterAuthenticationData struct {
	Type    string `json:"type"`
	Session string `json:"session"`
}

const LoginTypePassword = "m.login.password"
const LoginTypeToken = "m.login.token"

type LoginRequest struct {
	Type     string `json:"type"`
	Username string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	// ... and other parameters we don't care about
}

type ProfileDisplayNameRequest struct {
	DisplayName string `json:"displayname,omitempty"`
}

type ProfileAvatarUrlRequest struct {
	AvatarMxc string `json:"avatar_url,omitempty"`
}

type CreateRoomRequest struct {
	Visibility      string                 `json:"visibility,omitempty"`
	AliasLocalpart  string                 `json:"room_alias_name,omitempty"`
	InviteUserIds   []string               `json:"invite,flow,omitempty"`
	CreationContent map[string]interface{} `json:"creation_content,omitempty"`
	InitialState    []CreateRoomStateEvent `json:"initial_state,flow,omitempty"`
	Preset          string                 `json:"preset,omitempty"`
	IsDirect        bool                   `json:"is_direct"`
}

type CreateRoomStateEvent struct {
	Type     string      `json:"type"`
	StateKey string      `json:"state_key"`
	Content  interface{} `json:"content"`
}

type KickRequest struct {
	UserId string `json:"user_id"`
	Reason string `json:"reason,omitempty"`
}
