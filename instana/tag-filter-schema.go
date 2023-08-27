package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var tagFilterDiffSuppressFunc = func(k, old, new string, d *schema.ResourceData) bool {
	normalized, err := tagfilter.Normalize(new)
	if err == nil {
		return normalized == old
	}
	return old == new
}

var tagFilterStateFunc = func(val interface{}) string {
	normalized, err := tagfilter.Normalize(val.(string))
	if err == nil {
		return normalized
	}
	return val.(string)
}

var tagFilterValidateFunc = func(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if _, err := tagfilter.NewParser().Parse(v); err != nil {
		errs = append(errs, fmt.Errorf("%q is not a valid tag filter; %s", key, err))
	}

	return
}

var OptionalTagFilterExpressionSchema = &schema.Schema{
	Type:             schema.TypeString,
	Optional:         true,
	Description:      "The tag filter expression",
	DiffSuppressFunc: tagFilterDiffSuppressFunc,
	StateFunc:        tagFilterStateFunc,
	ValidateFunc:     tagFilterValidateFunc,
}

var RequiredTagFilterExpressionSchema = &schema.Schema{
	Type:             schema.TypeString,
	Required:         true,
	Description:      "The tag filter expression",
	DiffSuppressFunc: tagFilterDiffSuppressFunc,
	StateFunc:        tagFilterStateFunc,
	ValidateFunc:     tagFilterValidateFunc,
}
