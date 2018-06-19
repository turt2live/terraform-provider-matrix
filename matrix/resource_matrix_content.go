package matrix

import (
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
)

// TODO: Acceptance tests
// TODO: Register as matrix_content
// TODO: Require this for the avatar_mxc on users?

func resourceContent() *schema.Resource {
	return &schema.Resource{
		Create: resourceContentCreate,
		Read:   resourceContentRead,
		Update: resourceContentUpdate,
		Delete: resourceContentDelete,

		Schema: map[string]*schema.Schema{
			"mxc": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"media_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bytes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceContentCreate(d *schema.ResourceData, m interface{}) error {
	meta := m.(Metadata)

	mxcRaw := nilIfEmptyString(d.Get("mxc"))
	bytesRaw := nilIfEmptyString(d.Get("bytes"))
	typeRaw := nilIfEmptyString(d.Get("type"))
	fileNameRaw := nilIfEmptyString(d.Get("file_name"))

	if mxcRaw != nil && (bytesRaw != nil || typeRaw != nil || fileNameRaw != nil) {
		return fmt.Errorf("mxc uri cannot be provided alongside file information")
	}
	if mxcRaw == nil && bytesRaw == nil {
		return fmt.Errorf("either an mxc uri or content bytes must be supplied")
	}

	if mxcRaw != nil {
		mxc, origin, mediaId, err := stripMxc(mxcRaw.(string))
		if err != nil {
			return err
		}

		d.SetId(mxc)
		d.Set("mxc", mxc)
		d.Set("origin", origin)
		d.Set("media_id", mediaId)
	} else {
		// TODO: Upload content
		return fmt.Errorf("upload not yet implemented")
	}

	return resourceContentRead(d, meta)
}

func resourceContentRead(d *schema.ResourceData, m interface{}) error {
	// Nothing to do
	return nil
}

func resourceContentUpdate(d *schema.ResourceData, m interface{}) error {
	// There's nothing we can actually update
	// TODO: When the MXC changes, update the resource ID and other props
	// TODO: Are we able to detect changes to the bytes, etc and throw an error on that?
	// ... or do we just force a new resource?
	return nil
}

func resourceContentDelete(d *schema.ResourceData, m interface{}) error {
	// Content cannot be deleted in matrix (yet), so we just fake it
	return nil
}
