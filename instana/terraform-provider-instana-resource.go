package instana

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//TerraformProviderInstanaResource interface definition of a terraform resource of this terraform provider implementation
type TerraformProviderInstanaResource interface {
	Create(*schema.ResourceData, interface{}) error
	Read(d *schema.ResourceData, meta interface{}) error
	Update(d *schema.ResourceData, meta interface{}) error
	Delete(d *schema.ResourceData, meta interface{}) error

	GetSchema() map[string]*schema.Schema
	GetResourceName() string
}
