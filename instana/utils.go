package instana

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rs/xid"
)

//RandomID generates a random ID for a resource
func RandomID() string {
	xid := xid.New()
	return xid.String()
}

//ReadStringArrayParameterFromResource reads a string array parameter from a resource
func ReadStringArrayParameterFromResource(d *schema.ResourceData, key string) []string {

	if attr, ok := d.GetOk(key); ok {
		var array []string
		items := attr.([]interface{})
		for _, x := range items {
			item := x.(string)
			array = append(array, item)
		}

		return array
	}

	return nil
}

//ConvertSeverityFromInstanaAPIToTerraformRepresentation converts the integer representation of the Instana API to the string representation of the Terraform provider
func ConvertSeverityFromInstanaAPIToTerraformRepresentation(severity int) (string, error) {
	if severity == SeverityWarning.apiRepresentation {
		return SeverityWarning.terraformRepresentation, nil
	} else if severity == SeverityCritical.apiRepresentation {
		return SeverityCritical.terraformRepresentation, nil
	} else {
		return "INVALID", fmt.Errorf("%d is not a valid severity", severity)
	}
}

//ConvertSeverityFromTerraformToInstanaAPIRepresentation converts the string representation of the Terraform to the int representation of the Instana API provider
func ConvertSeverityFromTerraformToInstanaAPIRepresentation(severity string) (int, error) {
	if severity == SeverityWarning.terraformRepresentation {
		return SeverityWarning.apiRepresentation, nil
	} else if severity == SeverityCritical.terraformRepresentation {
		return SeverityCritical.apiRepresentation, nil
	} else {
		return -1, fmt.Errorf("%s is not a valid severity", severity)
	}
}
