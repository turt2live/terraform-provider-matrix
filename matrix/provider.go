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
		},

		ResourcesMap: map[string]*schema.Resource{
			"matrix_user": resourceUser(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Metadata{
		ClientApiUrl: d.Get("client_server_url").(string),
	}

	return config, nil
}
