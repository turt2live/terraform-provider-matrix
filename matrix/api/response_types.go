package api

import (
	"fmt"
)

type ErrorResponse struct {
	ErrorCode  string `json:"errcode"`
	Message    string `json:"error"`
	RawError   string
	StatusCode int
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("code=%s message=%s raw=%s status_code=%d", e.ErrorCode, e.Message, e.RawError, e.StatusCode)
}

type RegisterResponse struct {
	UserId      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	DeviceId    string `json:"device_id"`

	// home_server is deprecated and therefore not included
}

type LoginResponse struct {
	UserId      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	DeviceId    string `json:"device_id"`

	// home_server is deprecated and therefore not included
}

type ProfileResponse struct {
	DisplayName string `json:"displayname"`
	AvatarMxc   string `json:"avatar_url"`
}

type WhoAmIResponse struct {
	UserId string `json:"user_id"`
}

type AdminWhoisResponse struct {
	UserId string `json:"user_id"`
}

type UiAuthResponse struct {
	Session   string                 `json:"session"`
	Flows     []*UiAuthFlow          `json:"flows,flow"`
	Completed *[]string              `json:"completed,flow"`
	Params    map[string]interface{} `json:"params"`
}

type UiAuthFlow struct {
	Stages []string `json:"stages,flow"`
}

type ProfileUpdateResponse struct {
	// There isn't actually anything here
}

type ContentUploadResponse struct {
	ContentMxc string `json:"content_uri"`
}

type RoomIdResponse struct {
	RoomId string `json:"room_id"`
}

type EventIdResponse struct {
	EventId string `json:"event_id"`
}

type RoomDirectoryLookupResponse struct {
	RoomId  string `json:"room_id"`
	Servers []string `json:"servers,flow"`
}
