package restapi

import (
	"errors"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/utils"
)

func acValidateOpt(opt *string, err string) error {
	if opt == nil || utils.IsBlank(*opt) {
		return errors.New(err)
	}
	return nil
}

func acValidateOpts(opts map[string]*string) error {
	for k, v := range opts {
		if v == nil || utils.IsBlank(*v) {
			return errors.New(k)
		}
	}
	return nil
}

func acValidateList(opts []string, err string) error {
	if len(opts) == 0 {
		return errors.New(err)
	}
	return nil
}

func acValidateOpsGenieIntegration(apiKey *string, tags *string, region *OpsGenieRegionType) error {
	if apiKey == nil || utils.IsBlank(*apiKey) {
		return errors.New("API key is missing")
	}
	if tags == nil || utils.IsBlank(*tags) {
		return errors.New("tags are missing")
	}
	m := make(map[string]*string)
	m["API key is missing"] = apiKey
	m["Tags are missing"] = tags
	err := acValidateOpts(m)
	if err != nil {
		return err
	}
	if region == nil {
		return errors.New("region is missing")
	}
	if !IsSupportedOpsGenieRegionType(*region) {
		return fmt.Errorf("region %s is not valid", *region)
	}
	return nil
}
