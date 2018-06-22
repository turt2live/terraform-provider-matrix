package matrix

import (
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_server_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MATRIX_CLIENT_SERVER_URL", nil),
				Description: "The URL for your matrix homeserver. Eg: https://matrix.org",
			},
			"default_access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MATRIX_DEFAULT_ACCESS_TOKEN", ""),
				Description: "The default access token to use for miscellaneous requests (media uploads, etc)",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"matrix_user":    resourceUser(),
			"matrix_content": resourceContent(),
			"matrix_room":    resourceRoom(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Metadata{
		ClientApiUrl:       d.Get("client_server_url").(string),
		DefaultAccessToken: d.Get("default_access_token").(string),
	}

	return config, nil
}
