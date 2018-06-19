package matrix

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turt2live/terraform-provider-matrix/matrix/api"
	"log"
	"fmt"
	"errors"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_token": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"avatar_mxc": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	meta := m.(Metadata)

	usernameRaw := nilIfEmptyString(d.Get("username"))
	passwordRaw := nilIfEmptyString(d.Get("password"))
	accessTokenRaw := nilIfEmptyString(d.Get("access_token"))
	displayNameRaw := nilIfEmptyString(d.Get("display_name"))
	avatarMxcRaw := nilIfEmptyString(d.Get("avatar_mxc"))

	if passwordRaw == nil && accessTokenRaw == nil {
		return fmt.Errorf("either password or access_token must be supplied")
	}
	if passwordRaw != nil && accessTokenRaw != nil {
		return fmt.Errorf("both password and access_token cannot be supplied")
	}
	if passwordRaw != nil && usernameRaw == nil {
		return fmt.Errorf("username and password must be supplied")
	}

	if passwordRaw != nil {
		log.Printf("[DEBUG] User create: %s", usernameRaw.(string))
		response, err := api.DoRegister(meta.ClientApiUrl, usernameRaw.(string), passwordRaw.(string), "user")
		if err != nil {
			if r, ok := err.(*api.ErrorResponse); ok && r.ErrorCode == api.ErrCodeUserInUse {
				request := &api.LoginRequest{
					Type:     api.LoginTypePassword,
					Username: usernameRaw.(string),
					Password: passwordRaw.(string),
				}
				urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/login")
				response := &api.LoginResponse{}
				err2 := api.DoRequest("POST", urlStr, request, response, "")
				if err2 != nil {
					return fmt.Errorf("error logging in as user: %s", err)
				}

				d.SetId(response.UserId)
				d.Set("access_token", response.AccessToken)
			} else {
				return fmt.Errorf("error creating user: %s", err)
			}
		} else {
			d.SetId(response.UserId)
			d.Set("access_token", response.AccessToken)
		}
	} else {
		log.Printf("[DEBUG] User whoami")
		response := &api.WhoAmIResponse{}
		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/account/whoami")
		err := api.DoRequest("GET", urlStr, nil, response, accessTokenRaw.(string))
		if err != nil {
			return fmt.Errorf("error performing whoami: %s", err)
		}

		d.SetId(response.UserId)
	}

	if displayNameRaw != nil {
		resourceUserSetDisplayName(d, meta, displayNameRaw.(string))
	}

	if avatarMxcRaw != nil {
		resourceUserSetAvatarMxc(d, meta, avatarMxcRaw.(string))
	}

	return resourceUserRead(d, meta)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	meta := m.(Metadata)

	userId := d.Id()
	accessToken := d.Get("access_token").(string)

	log.Printf("[DEBUG] Getting user profile: %s", userId)
	urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/profile/", userId)
	response := &api.ProfileResponse{}
	err := api.DoRequest("GET", urlStr, nil, response, accessToken)
	if err != nil {
		if mtxErr, ok := err.(*api.ErrorResponse); ok && mtxErr.ErrorCode == api.ErrCodeUnknownToken {
			// Mark as deleted
			d.SetId("")
			d.Set("access_token", "")
			return nil
		}
		return fmt.Errorf("error getting user profile: %s", err)
	}

	d.Set("display_name", response.DisplayName)
	d.Set("avatar_mxc", response.AvatarMxc)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	// TODO: Change password
	// TODO: Change display name
	// TODO: Change avatar mxc
	return errors.New("not implemented")
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	// Users cannot be deleted in matrix, so we just say we deleted them
	return nil
}

func resourceUserSetDisplayName(d *schema.ResourceData, meta Metadata, newDisplayName string) error {
	accessToken := d.Get("access_token").(string)
	userId := d.Id()

	response := &api.ProfileUpdateResponse{}
	request := &api.ProfileDisplayNameRequest{DisplayName: newDisplayName}
	urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/profile/", userId, "/displayname")
	err := api.DoRequest("PUT", urlStr, request, response, accessToken)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserSetAvatarMxc(d *schema.ResourceData, meta Metadata, newAvatarMxc string) error {
	accessToken := d.Get("access_token").(string)
	userId := d.Id()

	response := &api.ProfileUpdateResponse{}
	request := &api.ProfileAvatarUrlRequest{AvatarMxc: newAvatarMxc}
	urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/profile/", userId, "/avatar_url")
	err := api.DoRequest("PUT", urlStr, request, response, accessToken)
	if err != nil {
		return err
	}

	return nil
}
