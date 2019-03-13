package instana

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Provider interface implementation of hashicorp terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{},
	}
}
