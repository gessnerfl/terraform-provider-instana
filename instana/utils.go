package instana

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rs/xid"
)

//RandomID generates a random ID for a resource
func RandomID() string {
	xid := xid.New()
	return xid.String()
}

//ReadStringArrayPtrFromResource reads a string array property from a resource
func ReadStringArrayPtrFromResource(d *schema.ResourceData, key string) []string {

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
