package matrix

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

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
}

func testAccAdminToken() string {
	return os.Getenv("MATRIX_ADMIN_ACCESS_TOKEN")
}
