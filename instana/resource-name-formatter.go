package instana

import (
	"strings"
)

//ResourceNameFormatter interface for the library to format resource name with a terraform managed string when configured
type ResourceNameFormatter interface {
	Format(name string) string
	UndoFormat(name string) string
}

//NewResourceNameFormatter creates a new formatter instance for the given prefix and suffix
func NewResourceNameFormatter(prefix string, suffix string) ResourceNameFormatter {
	return &terraformManagedResourceNameFormatter{
		prefix: prefix + " ",
		suffix: " " + suffix,
	}
}

//terraformManagedResourceNameFormatter implementation of ResourceNameFormatter which is used when terraform managed string should be appended to the name
type terraformManagedResourceNameFormatter struct {
	prefix string
	suffix string
}

func (formatter *terraformManagedResourceNameFormatter) Format(name string) string {
	return formatter.prefix + name + formatter.suffix
}

func (formatter *terraformManagedResourceNameFormatter) UndoFormat(name string) string {
	return strings.TrimPrefix(strings.TrimSuffix(name, formatter.suffix), formatter.prefix)
}
