package matrix

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"fmt"
	"github.com/turt2live/terraform-provider-matrix/matrix/api"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"regexp"
)

type testAccMatrixUser struct {
	Profile *api.ProfileResponse
	UserId  string
}

// TODO: Rename Basic test to be 'username password test'
// TODO: Test for when password but no username given
// TODO: Test for when password and access_token given
// TODO: Test for when no password and no access_token given
// TODO: Test for when username is taken (should login)
// TODO: Test for when password is wrong (cannot login)
// TODO: Test for when the access_token is invalid
// TODO: Test for when an access_token is given
// ... and probably other tests

// HACK: This test assumes the localpart (username) becomes the user ID for the user.
// From the spec: Matrix clients MUST NOT assume that localpart of the registered user_id matches the provided username.

var testAccCheckMatrixUserConfig_basic = fmt.Sprintf(`
resource "matrix_user" "foobar" {
	username = "foobar"
	password = "test1234"
}`)

func TestAccMatrixUser_Basic(t *testing.T) {
	var meta testAccMatrixUser

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if users get destroyed because they aren't
		//CheckDestroy: testAccCheckMatrixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMatrixUserConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixUserExists("matrix_user.foobar", &meta),
					testAccCheckMatrixUserAttributes(&meta, false, false),
					resource.TestMatchResourceAttr("matrix_user.foobar", "id", regexp.MustCompile("^@foobar:.*")),
					resource.TestCheckResourceAttr("matrix_user.foobar", "username", "foobar"),
					resource.TestCheckResourceAttr("matrix_user.foobar", "password", "test1234"),
					resource.TestMatchResourceAttr("matrix_user.foobar", "access_token", regexp.MustCompile(".+")),
					// we can't check the display name or avatar url because the homeserver might set it to something
				),
			},
		},
	})
}

func testAccCheckMatrixUserExists(n string, user *testAccMatrixUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		testUser := &testAccMatrixUser{}

		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/admin/whois/", rs.Primary.ID)
		response1 := &api.AdminWhoisResponse{}
		err := api.DoRequest("GET", urlStr, nil, response1, testAccAdminToken())
		if err != nil {
			return err
		}
		testUser.UserId = response1.UserId

		urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/profile/", rs.Primary.ID)
		response2 := &api.ProfileResponse{}
		err = api.DoRequest("GET", urlStr, nil, response2, testAccAdminToken())
		if err != nil {
			return err
		}
		testUser.Profile = response2

		*user = *testUser
		return nil
	}
}

func testAccCheckMatrixUserAttributes(user *testAccMatrixUser, checkDisplayName bool, checkAvatarUrl bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if strings.Index(user.UserId, "@foobar:") != 0 {
			return fmt.Errorf("bad user id: %s", user.UserId)
		}

		if checkDisplayName && user.Profile.DisplayName != "Baz" {
			return fmt.Errorf("bad display name: %s", user.Profile.DisplayName)
		}

		if checkAvatarUrl && user.Profile.AvatarMxc != "mxc://demo.site/abc123" {
			return fmt.Errorf("bad avatar mxc: %s", user.Profile.AvatarMxc)
		}

		return nil
	}
}
