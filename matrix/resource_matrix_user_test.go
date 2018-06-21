package matrix

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"fmt"
	"github.com/turt2live/terraform-provider-matrix/matrix/api"
	"github.com/hashicorp/terraform/terraform"
	"regexp"
)

type testAccMatrixUser struct {
	Profile *api.ProfileResponse
	UserId  string
}

// HACK: This test assumes the localpart (username) becomes the user ID for the user.
// From the spec: Matrix clients MUST NOT assume that localpart of the registered user_id matches the provided username.

var testAccMatrixUserConfig_usernamePassword = `
resource "matrix_user" "foobar" {
	username = "foobar"
	password = "test1234"
}`

func TestAccMatrixUser_UsernamePassword(t *testing.T) {
	var meta testAccMatrixUser

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if users get destroyed because they aren't
		//CheckDestroy: testAccCheckMatrixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMatrixUserConfig_usernamePassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixUserExists("matrix_user.foobar", &meta),
					testAccCheckMatrixUserIdMatches("matrix_user.foobar", &meta),
					testAccCheckMatrixUserAccessTokenWorks("matrix_user.foobar", &meta),
					resource.TestMatchResourceAttr("matrix_user.foobar", "id", regexp.MustCompile("^@foobar:.*")),
					resource.TestMatchResourceAttr("matrix_user.foobar", "access_token", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("matrix_user.foobar", "username", "foobar"),
					resource.TestCheckResourceAttr("matrix_user.foobar", "password", "test1234"),
					// we can't check the display name or avatar url because the homeserver might set it to something
				),
			},
		},
	})
}

var testAccMatrixUserConfig_usernamePasswordProfile = `
resource "matrix_user" "foobar" {
	username = "foobar"
	password = "test1234"
	display_name = "TEST_DISPLAY_NAME"
	avatar_mxc = "mxc://localhost/FakeAvatar"
}`

func TestAccMatrixUser_UsernamePasswordProfile(t *testing.T) {
	var meta testAccMatrixUser

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if users get destroyed because they aren't
		//CheckDestroy: testAccCheckMatrixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMatrixUserConfig_usernamePasswordProfile,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixUserExists("matrix_user.foobar", &meta),
					testAccCheckMatrixUserIdMatches("matrix_user.foobar", &meta),
					testAccCheckMatrixUserAccessTokenWorks("matrix_user.foobar", &meta),
					testAccCheckMatrixUserDisplayNameMatches("matrix_user.foobar", &meta),
					testAccCheckMatrixUserAvatarMxcMatches("matrix_user.foobar", &meta),
					resource.TestMatchResourceAttr("matrix_user.foobar", "id", regexp.MustCompile("^@foobar:.*")),
					resource.TestMatchResourceAttr("matrix_user.foobar", "access_token", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("matrix_user.foobar", "username", "foobar"),
					resource.TestCheckResourceAttr("matrix_user.foobar", "password", "test1234"),
					resource.TestCheckResourceAttr("matrix_user.foobar", "display_name", "TEST_DISPLAY_NAME"),
					resource.TestCheckResourceAttr("matrix_user.foobar", "avatar_mxc", "mxc://localhost/FakeAvatar"),
				),
			},
		},
	})
}

var testAccMatrixUserConfig_accessToken = `
resource "matrix_user" "foobar" {
	access_token = "%s"
}`

func TestAccMatrixUser_AccessToken(t *testing.T) {
	var meta testAccMatrixUser
	testUser := testAccCreateTestUser("test_user_access_token")
	conf := fmt.Sprintf(testAccMatrixUserConfig_accessToken, testUser.AccessToken)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if users get destroyed because they aren't
		//CheckDestroy: testAccCheckMatrixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixUserExists("matrix_user.foobar", &meta),
					testAccCheckMatrixUserIdMatches("matrix_user.foobar", &meta),
					//testAccCheckMatrixUserAccessTokenWorks("matrix_user.foobar", &meta),
					testAccCheckMatrixUserDisplayNameMatches("matrix_user.foobar", &meta),
					testAccCheckMatrixUserAvatarMxcMatches("matrix_user.foobar", &meta),
					resource.TestCheckResourceAttr("matrix_user.foobar", "id", testUser.UserId),
					resource.TestCheckResourceAttr("matrix_user.foobar", "access_token", testUser.AccessToken),
					resource.TestCheckResourceAttr("matrix_user.foobar", "display_name", testUser.DisplayName),
					resource.TestCheckResourceAttr("matrix_user.foobar", "avatar_mxc", testUser.AvatarMxc),
					resource.TestCheckNoResourceAttr("matrix_user.foobar", "username"),
					resource.TestCheckNoResourceAttr("matrix_user.foobar", "password"),
				),
			},
		},
	})
}

var testAccMatrixUserConfig_accessTokenProfile = `
resource "matrix_user" "foobar" {
	access_token = "%s"
	display_name = "%s"
	avatar_mxc = "%s"
}`

func TestAccMatrixUser_AccessTokenProfile(t *testing.T) {
	var meta testAccMatrixUser
	testUser := testAccCreateTestUser("test_user_access_token_profile")

	// We cheat and set the properties here to make sure they'll match the checks later on
	testUser.DisplayName = "TESTING1234"
	testUser.AvatarMxc = "mxc://localhost/SomeMediaID"

	conf := fmt.Sprintf(testAccMatrixUserConfig_accessTokenProfile, testUser.AccessToken, testUser.DisplayName, testUser.AvatarMxc)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if users get destroyed because they aren't
		//CheckDestroy: testAccCheckMatrixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixUserExists("matrix_user.foobar", &meta),
					testAccCheckMatrixUserIdMatches("matrix_user.foobar", &meta),
					//testAccCheckMatrixUserAccessTokenWorks("matrix_user.foobar", &meta),
					testAccCheckMatrixUserDisplayNameMatches("matrix_user.foobar", &meta),
					testAccCheckMatrixUserAvatarMxcMatches("matrix_user.foobar", &meta),
					resource.TestCheckResourceAttr("matrix_user.foobar", "id", testUser.UserId),
					resource.TestCheckResourceAttr("matrix_user.foobar", "access_token", testUser.AccessToken),
					resource.TestCheckResourceAttr("matrix_user.foobar", "display_name", testUser.DisplayName),
					resource.TestCheckResourceAttr("matrix_user.foobar", "avatar_mxc", testUser.AvatarMxc),
					resource.TestCheckNoResourceAttr("matrix_user.foobar", "username"),
					resource.TestCheckNoResourceAttr("matrix_user.foobar", "password"),
				),
			},
		},
	})
}

var testAccMatrixUserConfig_updateProfile = `
resource "matrix_user" "foobar" {
	access_token = "%s"
	display_name = "%s"
	avatar_mxc = "%s"
}`

func TestAccMatrixUser_UpdateProfile(t *testing.T) {
	var meta testAccMatrixUser
	originalUser := testAccCreateTestUser("test_user_update_profile")

	// We cheat and set the properties here to make sure they'll match the checks later on
	originalUser.DisplayName = "TESTING1234"
	originalUser.AvatarMxc = "mxc://localhost/SomeMediaID"

	updatedUser := &test_MatrixUser{
		UserId:      originalUser.UserId,
		AvatarMxc:   "mxc://localhost/SomeOtherMediaId",
		DisplayName: "New Display Name",
		AccessToken: originalUser.AccessToken,
		Localpart:   originalUser.Localpart,
		Password:    originalUser.Password,
	}

	confPart1 := fmt.Sprintf(testAccMatrixUserConfig_updateProfile, originalUser.AccessToken, originalUser.DisplayName, originalUser.AvatarMxc)
	confPart2 := fmt.Sprintf(testAccMatrixUserConfig_updateProfile, updatedUser.AccessToken, updatedUser.DisplayName, updatedUser.AvatarMxc)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// We don't check if users get destroyed because they aren't
		//CheckDestroy: testAccCheckMatrixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: confPart1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixUserExists("matrix_user.foobar", &meta),
					testAccCheckMatrixUserIdMatches("matrix_user.foobar", &meta),
					//testAccCheckMatrixUserAccessTokenWorks("matrix_user.foobar", &meta),
					testAccCheckMatrixUserDisplayNameMatches("matrix_user.foobar", &meta),
					testAccCheckMatrixUserAvatarMxcMatches("matrix_user.foobar", &meta),
					resource.TestCheckResourceAttr("matrix_user.foobar", "id", originalUser.UserId),
					resource.TestCheckResourceAttr("matrix_user.foobar", "access_token", originalUser.AccessToken),
					resource.TestCheckResourceAttr("matrix_user.foobar", "display_name", originalUser.DisplayName),
					resource.TestCheckResourceAttr("matrix_user.foobar", "avatar_mxc", originalUser.AvatarMxc),
					resource.TestCheckNoResourceAttr("matrix_user.foobar", "username"),
					resource.TestCheckNoResourceAttr("matrix_user.foobar", "password"),
				),
			},
			{
				Config: confPart2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatrixUserExists("matrix_user.foobar", &meta),
					testAccCheckMatrixUserIdMatches("matrix_user.foobar", &meta),
					//testAccCheckMatrixUserAccessTokenWorks("matrix_user.foobar", &meta),
					testAccCheckMatrixUserDisplayNameMatches("matrix_user.foobar", &meta),
					testAccCheckMatrixUserAvatarMxcMatches("matrix_user.foobar", &meta),
					resource.TestCheckResourceAttr("matrix_user.foobar", "id", updatedUser.UserId),
					resource.TestCheckResourceAttr("matrix_user.foobar", "access_token", updatedUser.AccessToken),
					resource.TestCheckResourceAttr("matrix_user.foobar", "display_name", updatedUser.DisplayName),
					resource.TestCheckResourceAttr("matrix_user.foobar", "avatar_mxc", updatedUser.AvatarMxc),
					resource.TestCheckNoResourceAttr("matrix_user.foobar", "username"),
					resource.TestCheckNoResourceAttr("matrix_user.foobar", "password"),
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

		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/admin/whois/", rs.Primary.ID)
		response1 := &api.AdminWhoisResponse{}
		err := api.DoRequest("GET", urlStr, nil, response1, testAccAdminToken())
		if err != nil {
			return err
		}

		urlStr = api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/profile/", rs.Primary.ID)
		response2 := &api.ProfileResponse{}
		err = api.DoRequest("GET", urlStr, nil, response2, testAccAdminToken())
		if err != nil {
			return err
		}

		*user = testAccMatrixUser{
			UserId:  response1.UserId,
			Profile: response2,
		}
		return nil
	}
}

func testAccCheckMatrixUserIdMatches(n string, user *testAccMatrixUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		if rs.Primary.ID != user.UserId {
			return fmt.Errorf("user id doesn't match. expected: %s  got: %s", user.UserId, rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckMatrixUserDisplayNameMatches(n string, user *testAccMatrixUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		displayNameRaw := nilIfEmptyString(rs.Primary.Attributes["display_name"])
		if displayNameRaw != user.Profile.DisplayName {
			return fmt.Errorf("display name doesn't match. exepcted: %s  got: %s", user.Profile.DisplayName, displayNameRaw)
		}

		return nil
	}
}

func testAccCheckMatrixUserAvatarMxcMatches(n string, user *testAccMatrixUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		avatarMxcRaw := nilIfEmptyString(rs.Primary.Attributes["avatar_mxc"])
		if avatarMxcRaw != user.Profile.AvatarMxc {
			return fmt.Errorf("display name doesn't match. exepcted: %s  got: %s", user.Profile.AvatarMxc, avatarMxcRaw)
		}

		return nil
	}
}

func testAccCheckMatrixUserAccessTokenWorks(n string, user *testAccMatrixUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		meta := testAccProvider.Meta().(Metadata)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("record id not set")
		}

		accessTokenRaw := nilIfEmptyString(rs.Primary.Attributes["access_token"])

		response := &api.WhoAmIResponse{}
		urlStr := api.MakeUrl(meta.ClientApiUrl, "/_matrix/client/r0/account/whoami")
		err := api.DoRequest("GET", urlStr, nil, response, accessTokenRaw.(string))
		if err != nil {
			return fmt.Errorf("error performing whoami: %s", err)
		}

		if response.UserId != rs.Primary.ID {
			return fmt.Errorf("whoami succeeded, although the user id does not match. expected: %s  got: %s", rs.Primary.ID, response.UserId)
		}

		return nil
	}
}
