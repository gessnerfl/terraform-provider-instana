package instana

import (
	"github.com/rs/xid"
)

//RandomID generates a random ID for a resource
func RandomID() string {
	xid := xid.New()
	return xid.String()
}
