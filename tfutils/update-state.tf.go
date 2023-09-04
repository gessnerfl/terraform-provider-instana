package tfutils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func UpdateState(d *schema.ResourceData, data map[string]interface{}) error {
	for k, v := range data {
		err := d.Set(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
