package instana

import (
	"strings"
)

//TerraformManagedResourceNameSuffix the suffix which is appended to a name string
const TerraformManagedResourceNameSuffix = " (TF managed)"

//ResourceNameFormatter interface for the library to format resource name with a terraform managed string when configured
type ResourceNameFormatter interface {
	Format(name string) string
	UndoFormat(name string) string
}

//NewResourceNameFormatter creates a new formatter instance depending on if terraform managed string append is requested or not
func NewResourceNameFormatter(appendString bool) ResourceNameFormatter {
	if appendString {
		return &terraformManagedResourceNameFormatter{}
	}
	return &noopResourceNameFormatter{}
}

//noopResourceNameFormatter implementation of ResourceNameFormatter which is used when no terraform managed string should be appended to the name
type noopResourceNameFormatter struct{}

func (d *noopResourceNameFormatter) Format(name string) string {
	return name
}

func (d *noopResourceNameFormatter) UndoFormat(name string) string {
	return name
}

//terraformManagedResourceNameFormatter implementation of ResourceNameFormatter which is used when terraform managed string should be appended to the name
type terraformManagedResourceNameFormatter struct{}

func (d *terraformManagedResourceNameFormatter) Format(name string) string {
	return name + TerraformManagedResourceNameSuffix
}

func (d *terraformManagedResourceNameFormatter) UndoFormat(name string) string {
	return strings.TrimSuffix(name, TerraformManagedResourceNameSuffix)
}
