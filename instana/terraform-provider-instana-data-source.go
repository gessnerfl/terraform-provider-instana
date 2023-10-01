package instana

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// DataSource interface definition of a Terraform DataSource implementation in this provider
type DataSource interface {
	CreateResource() *schema.Resource
}
