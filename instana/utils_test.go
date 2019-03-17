package instana_test

import (
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"testing"
)

func TestRandomID(t *testing.T) {
	id := RandomID()

	if len(id) == 0 {
		t.Fatal("Expected to get a new id generated")
	}
}
