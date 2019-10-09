package instana

//ResourceStringFormatter interface for the library to format resource names and/or description with a terraform managed string when configured
type ResourceStringFormatter interface {
	FormatName(name string) string
	FormatDescription(description string) string
}

//NewResourceStringFormatter creates a new formatter instance depending on if terraform managed string append is requested or not
func NewResourceStringFormatter(appendString bool) ResourceStringFormatter {
	if appendString {
		return &terraformManagedResourceStringFormatter{}
	}
	return &noopResourceStringFormatter{}
}

//noopResourceStringFormatter implementation of ResourceStringFormatter which is used when no terraform managed string should be appended
type noopResourceStringFormatter struct{}

func (d *noopResourceStringFormatter) FormatName(name string) string {
	return name
}

func (d *noopResourceStringFormatter) FormatDescription(description string) string {
	return description
}

//terraformManagedResourceStringFormatter implementation of ResourceStringFormatter which is used when terraform managed string should be appended
type terraformManagedResourceStringFormatter struct{}

func (d *terraformManagedResourceStringFormatter) FormatName(name string) string {
	return name + " (tf managed)"
}

func (d *terraformManagedResourceStringFormatter) FormatDescription(description string) string {
	return description + "\n\n--\nterraform managed"
}
