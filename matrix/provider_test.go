package matrix

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
	"github.com/turt2live/terraform-provider-matrix/matrix/api"
)

type test_MatrixUser struct {
	Localpart   string
	Password    string
	AccessToken string
	UserId      string
	DisplayName string
	AvatarMxc   string
}

var test_MatrixUser_users = make(map[string]*test_MatrixUser)
var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"matrix": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("MATRIX_CLIENT_SERVER_URL"); v == "" {
		t.Fatal("MATRIX_CLIENT_SERVER_URL must be set for acceptance tests")
	}
	if v := os.Getenv("MATRIX_ADMIN_ACCESS_TOKEN"); v == "" {
		t.Fatal("MATRIX_ADMIN_ACCESS_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("MATRIX_DEFAULT_ACCESS_TOKEN"); v == "" {
		t.Fatal("MATRIX_DEFAULT_ACCESS_TOKEN must be set for acceptance tests")
	}
}

func testAccTestDataDir() string {
	return os.Getenv("MATRIX_TEST_DATA_DIR")
}

func testAccClientServerUrl() string {
	return os.Getenv("MATRIX_CLIENT_SERVER_URL")
}

func testAccAdminToken() string {
	return os.Getenv("MATRIX_ADMIN_ACCESS_TOKEN")
}

func testAccCreateTestUser(localpart string) (*test_MatrixUser) {
	existing := test_MatrixUser_users[localpart]
	if existing != nil {
		return existing
	}

	csApiUrl := testAccProvider.Meta().(Metadata).ClientApiUrl
	password := "test1234"
	displayName := "!!TEST USER!!"
	avatarMxc := "mxc://domain.com/SomeAvatarUrl"

	r, e := api.DoRegister(csApiUrl, localpart, password, "user")
	if e != nil {
		panic(e)
	}

	response := &api.ProfileUpdateResponse{}
	nameRequest := &api.ProfileDisplayNameRequest{DisplayName: displayName}
	urlStr := api.MakeUrl(csApiUrl, "/_matrix/client/r0/profile/", r.UserId, "/displayname")
	e = api.DoRequest("PUT", urlStr, nameRequest, response, r.AccessToken)
	if e != nil {
		panic(e)
	}

	avatarRequest := &api.ProfileAvatarUrlRequest{AvatarMxc: avatarMxc}
	urlStr = api.MakeUrl(csApiUrl, "/_matrix/client/r0/profile/", r.UserId, "/avatar_url")
	e = api.DoRequest("PUT", urlStr, avatarRequest, response, r.AccessToken)
	if e != nil {
		panic(e)
	}

	existing = &test_MatrixUser{
		Localpart:   localpart,
		Password:    password,
		AccessToken: r.AccessToken,
		UserId:      r.UserId,
		DisplayName: displayName,
		AvatarMxc:   avatarMxc,
	}

	test_MatrixUser_users[localpart] = existing
	return existing
}
