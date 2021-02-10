package instana_test

import (
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const contentType = "Content-Type"

var testProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}
