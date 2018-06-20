package matrix

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"fmt"
	"github.com/turt2live/terraform-provider-matrix/matrix/api"
	"github.com/hashicorp/terraform/terraform"
	"regexp"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path"
)

type testAccMatrixContentUpload struct {
	Mxc      string
	Origin   string
	MediaId  string
	Content  []byte
	FilePath string
	FileName string
	FileType string
}

func testAccCreateMatrixContent(content []byte, mime string, fileName string) (*testAccMatrixContentUpload) {
	response, err := api.UploadFile(testAccClientServerUrl(), content, fileName, mime, testAccAdminToken())
	if err != nil {
		panic(err)
	}

	mxc, origin, mediaId, err := stripMxc(response.ContentMxc)
	if err != nil {
		panic(err)
	}

	return &testAccMatrixContentUpload{
		Mxc:      mxc,
		Origin:   origin,
		MediaId:  mediaId,
		Content:  content,
		FilePath: "",
		FileName: fileName,
		FileType: mime,
	}
}

var testAccMatrixContentConfig_existingContent = `
resource "matrix_content" "foobar" {
	origin = "%s"
	media_id = "%s"
}`

func TestAccMatrixContent_ExistingContent(t *testing.T) {
	upload := testAccCreateMatrixContent([]byte("hello world"), "text/plain", "hello.txt")
	conf := fmt.Sprintf(testAccMatrixContentConfig_existingContent, upload.Origin, upload.MediaId)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if content get destroyed because it isn't
		//CheckDestroy: testAccCheckMatrixContentDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixContentExists("matrix_content.foobar"),
					testAccCheckMatrixContentMatchesUpload("matrix_content.foobar", upload),
					testAccCheckMatrixContentIdMatchesProperties("matrix_content.foobar"),
					resource.TestMatchResourceAttr("matrix_content.foobar", "id", regexp.MustCompile("^mxc://[a-zA-Z0-9.:\\-_]+/[a-zA-Z0-9]+")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "origin", regexp.MustCompile("^[a-zA-Z0-9.:\\-_]+$")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "media_id", regexp.MustCompile("^[a-zA-Z0-9]+$")),
				),
			},
		},
	})
}

var testAccMatrixContentConfig_plainText = `
resource "matrix_content" "foobar" {
	file_path = "%s"
	file_name = "%s"
	file_type = "%s"
}`

func TestAccMatrixContent_PlainTextUpload(t *testing.T) {
	upload := &testAccMatrixContentUpload{
		FilePath: path.Join(testAccTestDataDir(), ".test_data/words.txt"),
		FileName: "hello.txt",
		FileType: "text/plain",
	}
	conf := fmt.Sprintf(testAccMatrixContentConfig_plainText, upload.FilePath, upload.FileName, upload.FileType)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if content get destroyed because it isn't
		//CheckDestroy: testAccCheckMatrixContentDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixContentExists("matrix_content.foobar"),
					//testAccCheckMatrixContentMatchesUpload("matrix_content.foobar", upload),
					testAccCheckMatrixContentMatchesFile("matrix_content.foobar", upload),
					testAccCheckMatrixContentIdMatchesProperties("matrix_content.foobar"),
					resource.TestMatchResourceAttr("matrix_content.foobar", "id", regexp.MustCompile("^mxc://[a-zA-Z0-9.:\\-_]+/[a-zA-Z0-9]+")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "origin", regexp.MustCompile("^[a-zA-Z0-9.:\\-_]+$")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "media_id", regexp.MustCompile("^[a-zA-Z0-9]+$")),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_path", upload.FilePath),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_name", upload.FileName),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_type", upload.FileType),
				),
			},
		},
	})
}

var testAccMatrixContentConfig_binary = `
resource "matrix_content" "foobar" {
	file_path = "%s"
	file_name = "%s"
	file_type = "%s"
}`

func TestAccMatrixContent_BinaryUpload(t *testing.T) {
	upload := &testAccMatrixContentUpload{
		FilePath: path.Join(testAccTestDataDir(), ".test_data/deadbeef.bin"),
		FileName: "beef.bin",
		FileType: "application/octet-stream",
	}
	conf := fmt.Sprintf(testAccMatrixContentConfig_plainText, upload.FilePath, upload.FileName, upload.FileType)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if content get destroyed because it isn't
		//CheckDestroy: testAccCheckMatrixContentDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixContentExists("matrix_content.foobar"),
					//testAccCheckMatrixContentMatchesUpload("matrix_content.foobar", upload),
					testAccCheckMatrixContentMatchesFile("matrix_content.foobar", upload),
					testAccCheckMatrixContentIdMatchesProperties("matrix_content.foobar"),
					resource.TestMatchResourceAttr("matrix_content.foobar", "id", regexp.MustCompile("^mxc://[a-zA-Z0-9.:\\-_]+/[a-zA-Z0-9]+")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "origin", regexp.MustCompile("^[a-zA-Z0-9.:\\-_]+$")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "media_id", regexp.MustCompile("^[a-zA-Z0-9]+$")),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_path", upload.FilePath),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_name", upload.FileName),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_type", upload.FileType),
				),
			},
		},
	})
}

var testAccMatrixContentConfig_uploadNoName = `
resource "matrix_content" "foobar" {
	file_path = "%s"
	file_type = "%s"
}`

func TestAccMatrixContent_UploadNoName(t *testing.T) {
	upload := &testAccMatrixContentUpload{
		FilePath: path.Join(testAccTestDataDir(), ".test_data/deadbeef.bin"),
		FileName: "", // expected name
		FileType: "application/octet-stream",
	}
	conf := fmt.Sprintf(testAccMatrixContentConfig_uploadNoName, upload.FilePath, upload.FileType)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if content get destroyed because it isn't
		//CheckDestroy: testAccCheckMatrixContentDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixContentExists("matrix_content.foobar"),
					//testAccCheckMatrixContentMatchesUpload("matrix_content.foobar", upload),
					testAccCheckMatrixContentMatchesFile("matrix_content.foobar", upload),
					testAccCheckMatrixContentIdMatchesProperties("matrix_content.foobar"),
					resource.TestMatchResourceAttr("matrix_content.foobar", "id", regexp.MustCompile("^mxc://[a-zA-Z0-9.:\\-_]+/[a-zA-Z0-9]+")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "origin", regexp.MustCompile("^[a-zA-Z0-9.:\\-_]+$")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "media_id", regexp.MustCompile("^[a-zA-Z0-9]+$")),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_path", upload.FilePath),
					resource.TestCheckNoResourceAttr("matrix_content.foobar", "file_name"),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_type", upload.FileType),
				),
			},
		},
	})
}

var testAccMatrixContentConfig_uploadNoType = `
resource "matrix_content" "foobar" {
	file_path = "%s"
	file_name = "%s"
}`

func TestAccMatrixContent_UploadNoType(t *testing.T) {
	upload := &testAccMatrixContentUpload{
		FilePath: path.Join(testAccTestDataDir(), ".test_data/deadbeef.bin"),
		FileName: "beef.bin",
		FileType: "application/octet-stream", // expected type
	}
	conf := fmt.Sprintf(testAccMatrixContentConfig_uploadNoType, upload.FilePath, upload.FileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if content get destroyed because it isn't
		//CheckDestroy: testAccCheckMatrixContentDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixContentExists("matrix_content.foobar"),
					//testAccCheckMatrixContentMatchesUpload("matrix_content.foobar", upload),
					testAccCheckMatrixContentMatchesFile("matrix_content.foobar", upload),
					testAccCheckMatrixContentIdMatchesProperties("matrix_content.foobar"),
					resource.TestMatchResourceAttr("matrix_content.foobar", "id", regexp.MustCompile("^mxc://[a-zA-Z0-9.:\\-_]+/[a-zA-Z0-9]+")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "origin", regexp.MustCompile("^[a-zA-Z0-9.:\\-_]+$")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "media_id", regexp.MustCompile("^[a-zA-Z0-9]+$")),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_path", upload.FilePath),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_name", upload.FileName),
					resource.TestCheckNoResourceAttr("matrix_content.foobar", "file_type"),
				),
			},
		},
	})
}

var testAccMatrixContentConfig_uploadNoNameOrType = `
resource "matrix_content" "foobar" {
	file_path = "%s"
}`

func TestAccMatrixContent_UploadNoNameOrType(t *testing.T) {
	upload := &testAccMatrixContentUpload{
		FilePath: path.Join(testAccTestDataDir(), ".test_data/deadbeef.bin"),
		FileName: "",                         // expected name
		FileType: "application/octet-stream", // expected type
	}
	conf := fmt.Sprintf(testAccMatrixContentConfig_uploadNoNameOrType, upload.FilePath)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if content get destroyed because it isn't
		//CheckDestroy: testAccCheckMatrixContentDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixContentExists("matrix_content.foobar"),
					//testAccCheckMatrixContentMatchesUpload("matrix_content.foobar", upload),
					testAccCheckMatrixContentMatchesFile("matrix_content.foobar", upload),
					testAccCheckMatrixContentIdMatchesProperties("matrix_content.foobar"),
					resource.TestMatchResourceAttr("matrix_content.foobar", "id", regexp.MustCompile("^mxc://[a-zA-Z0-9.:\\-_]+/[a-zA-Z0-9]+")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "origin", regexp.MustCompile("^[a-zA-Z0-9.:\\-_]+$")),
					resource.TestMatchResourceAttr("matrix_content.foobar", "media_id", regexp.MustCompile("^[a-zA-Z0-9]+$")),
					resource.TestCheckResourceAttr("matrix_content.foobar", "file_path", upload.FilePath),
					resource.TestCheckNoResourceAttr("matrix_content.foobar", "file_name"),
					resource.TestCheckNoResourceAttr("matrix_content.foobar", "file_type"),
				),
			},
		},
	})
}

func testAccCheckMatrixContentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		origin := rs.Primary.Attributes["origin"]
		mediaId := rs.Primary.Attributes["media_id"]

		stream, _, err := api.DownloadFile(meta.ClientApiUrl, origin, mediaId)
		if stream != nil {
			defer (*stream).Close()
			io.Copy(ioutil.Discard, *stream)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckMatrixContentMatchesFile(n string, upload *testAccMatrixContentUpload) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		origin := rs.Primary.Attributes["origin"]
		mediaId := rs.Primary.Attributes["media_id"]

		download, headers, err := api.DownloadFile(meta.ClientApiUrl, origin, mediaId)
		contents := make([]byte, 0)
		if download != nil {
			defer (*download).Close()
			contents, err = ioutil.ReadAll(*download)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}

		// Ensure the Content is populated
		if upload.FilePath != "" {
			f, err := os.Open(upload.FilePath)
			if err != nil {
				return err
			}
			defer f.Close()

			b, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			upload.Content = b
		}

		// Compare bytes
		if len(contents) != len(upload.Content) {
			return fmt.Errorf("content length mismatch. expected: %d  got: %d", len(upload.Content), len(contents))
		}
		for i := range contents {
			d := contents[i]
			e := upload.Content[i]
			if d != e {
				return fmt.Errorf("byte mismatch at index %d/%d. expected: %b  got: %b", i, len(contents), e, d)
			}
		}

		// Compare content type
		contentType := headers.Get("Content-Type")
		if contentType != upload.FileType {
			return fmt.Errorf("content type mismatch. expected: %s  got: %s", upload.FileType, contentType)
		}

		// Compare file name
		contentDisposition := headers.Get("content-disposition")
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			if err.Error() != "mime: no media type" {
				return err
			}

			params = map[string]string{"filename": ""}
		}
		fileName := params["filename"]
		if fileName != upload.FileName {
			return fmt.Errorf("file name mismatch. expected: %s  got: %s", upload.FileName, fileName)
		}

		return nil
	}
}

func testAccCheckMatrixContentMatchesUpload(n string, uploaded *testAccMatrixContentUpload) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		mxc := rs.Primary.ID
		origin := rs.Primary.Attributes["origin"]
		mediaId := rs.Primary.Attributes["media_id"]

		if mxc != uploaded.Mxc {
			return fmt.Errorf("mxc does not match. expected: %s  got: %s", uploaded.Mxc, mxc)
		}
		if origin != uploaded.Origin {
			return fmt.Errorf("origin does not match. expected: %s  got: %s", uploaded.Origin, origin)
		}
		if mediaId != uploaded.MediaId {
			return fmt.Errorf("media_id does not match. expected: %s  got: %s", uploaded.MediaId, mediaId)
		}

		return nil
	}
}

func testAccCheckMatrixContentIdMatchesProperties(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		mxc := rs.Primary.ID
		origin := rs.Primary.Attributes["origin"]
		mediaId := rs.Primary.Attributes["media_id"]

		calcMxc := fmt.Sprintf("mxc://%s/%s", origin, mediaId)
		if calcMxc != mxc {
			return fmt.Errorf("id and calculated mxc are different. expected: %s  got: %s", calcMxc, mxc)
		}

		return nil
	}
}
