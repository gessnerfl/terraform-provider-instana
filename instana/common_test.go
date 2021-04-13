package instana_test

import (
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const contentType = "Content-Type"

var testProviders = map[string]*schema.Provider{
	"instana": Provider(),
}

func getZeroBasedCallCount(httpServer testutils.TestHTTPServer, method string, path string) int {
	count := httpServer.GetCallCount(method, path)
	if count == 0 {
		return count
	}
	return count - 1
}
