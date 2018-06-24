package api

import (
	"net/http"
	"encoding/json"
	"errors"
	"log"
)

const AuthTypeDummy = "m.login.dummy"

func DoRegister(csApiUrl string, username string, password string, kind string) (*RegisterResponse, error) {
	qs := map[string]string{"kind": kind}
	urlStr := MakeUrlQueryString(qs, csApiUrl, "/_matrix/client/r0/register")

	// First we do a request to get the flows we can use
	log.Println("[DEBUG] Getting registration flows")
	request := &RegisterRequest{}
	state, _, err := doUiAuthRegisterRequest(urlStr, request)
	if err != nil {
		return nil, err
	}

	// Now that we have the process started, make sure we can actually follow one of these methods
	hasDummyStage := false
	for _, flow := range state.Flows {
		if len(flow.Stages) == 1 && flow.Stages[0] == AuthTypeDummy {
			hasDummyStage = true
			break
		}
	}
	if !hasDummyStage {
		return nil, errors.New("no dummy auth stage")
	}

	// We have a dummy stage, so we can expect to be able to register now
	log.Println("[DEBUG] Using dummy registration flow to register user")
	request = &RegisterRequest{
		Authentication: &RegisterAuthenticationData{
			Type:    AuthTypeDummy,
			Session: state.Session,
		},
		Username: username,
		Password: password,
	}
	_, response, err := doUiAuthRegisterRequest(urlStr, request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("ui auth failed: expected response but got login flow")
	}

	return response, nil
}

func doUiAuthRegisterRequest(urlStr string, request *RegisterRequest) (*UiAuthResponse, *RegisterResponse, error) {
	response := &RegisterResponse{}
	err := DoRequest("POST", urlStr, request, response, "")
	if err != nil {
		if r, ok := err.(*ErrorResponse); ok {
			if r.StatusCode == http.StatusUnauthorized {
				authState := &UiAuthResponse{}
				err2 := json.Unmarshal([]byte(r.RawError), authState)
				if err2 != nil {
					return nil, nil, err2
				}

				return authState, nil, nil
			}
		}

		return nil, nil, err
	}

	return nil, response, nil
}
