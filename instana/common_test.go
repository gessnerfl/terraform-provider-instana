package instana_test

import (
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const contentType = "Content-Type"

var testProviders = map[string]*schema.Provider{
	"instana": Provider(),
}
