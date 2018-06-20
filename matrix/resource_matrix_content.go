package matrix

import (
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
	"github.com/turt2live/terraform-provider-matrix/matrix/api"
	"io"
	"io/ioutil"
	"os"
)

// TODO: Acceptance tests: OnChange

func resourceContent() *schema.Resource {
	return &schema.Resource{
		Exists: resourceContentExists,
		Create: resourceContentCreate,
		Read:   resourceContentRead,
		Update: resourceContentUpdate,
		Delete: resourceContentDelete,

		Schema: map[string]*schema.Schema{
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"media_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_type": {
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

	originRaw := nilIfEmptyString(d.Get("origin"))
	mediaIdRaw := nilIfEmptyString(d.Get("media_id"))
	filePathRaw := nilIfEmptyString(d.Get("file_path"))
	fileTypeRaw := nilIfEmptyString(d.Get("file_type"))
	fileNameRaw := nilIfEmptyString(d.Get("file_name"))

	if (originRaw != nil && mediaIdRaw == nil) || (originRaw == nil && mediaIdRaw != nil) {
		return fmt.Errorf("both the media_id and origin must be supplied")
	}

	var mxcRaw interface{}
	if originRaw != nil {
		mxcRaw = fmt.Sprintf("mxc://%s/%s", originRaw, mediaIdRaw)
	}

	if mxcRaw != nil && (filePathRaw != nil || fileTypeRaw != nil || fileNameRaw != nil) {
		return fmt.Errorf("origin and media_id cannot be provided alongside file information")
	}
	if mxcRaw == nil && filePathRaw == nil {
		return fmt.Errorf("file_path must be supplied or an origin with media_id")
	}

	if mxcRaw != nil {
		mxc, origin, mediaId, err := stripMxc(mxcRaw.(string))
		if err != nil {
			return err
		}

		if origin != originRaw {
			return fmt.Errorf("origin mismatch while creating object. expected: '%s'  got: '%s'", originRaw, origin)
		}
		if mediaId != mediaIdRaw {
			return fmt.Errorf("media_id mismatch while creating object. expected: '%s'  got: '%s'", mediaIdRaw, mediaId)
		}

		d.SetId(mxc)
		d.Set("origin", origin)
		d.Set("media_id", mediaId)
	} else {
		if meta.DefaultAccessToken == "" {
			return fmt.Errorf("a default access token is required to upload content")
		}

		f, err := os.Open(filePathRaw.(string))
		if err != nil {
			return fmt.Errorf("error opening file: %s", err)
		}
		defer f.Close()

		contentBytes, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err)
		}

		fileName := ""
		if fileNameRaw != nil {
			fileName = fileNameRaw.(string)
		}

		contentType := "application/octet-stream"
		if fileTypeRaw != nil {
			contentType = fileTypeRaw.(string)
		}

		result, err := api.UploadFile(meta.ClientApiUrl, contentBytes, fileName, contentType, meta.DefaultAccessToken)
		if err != nil {
			return fmt.Errorf("error uploading content: %s", err)
		}

		mxc, origin, mediaId, err := stripMxc(result.ContentMxc)
		if err != nil {
			return err
		}

		d.SetId(mxc)
		d.Set("origin", origin)
		d.Set("media_id", mediaId)
	}

	return resourceContentRead(d, meta)
}

func resourceContentExists(d *schema.ResourceData, m interface{}) (bool, error) {
	meta := m.(Metadata)

	origin := d.Get("origin").(string)
	mediaId := d.Get("media_id").(string)

	stream, _, err := api.DownloadFile(meta.ClientApiUrl, origin, mediaId)
	if stream != nil {
		defer (*stream).Close()
		io.Copy(ioutil.Discard, *stream)
	}
	if err != nil {
		return false, nil
	}

	return true, nil
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

	/*
		// Enable partial state mode
        d.Partial(true)

        if d.HasChange("address") {
                // Try updating the address
                if err := updateAddress(d, m); err != nil {
                        return err
                }

                d.SetPartial("address")
        }

        // If we were to return here, before disabling partial mode below,
        // then only the "address" field would be saved.

        // We succeeded, disable partial mode. This causes Terraform to save
        // all fields again.
        d.Partial(false)

        return nil
	 */

	return nil
}

func resourceContentDelete(d *schema.ResourceData, m interface{}) error {
	// Content cannot be deleted in matrix (yet), so we just fake it
	return nil
}
