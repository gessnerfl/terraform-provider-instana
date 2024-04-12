package restapi

type CustomPayloadFieldsAware interface {
	GetCustomerPayloadFields() []CustomPayloadField[any]
	SetCustomerPayloadFields([]CustomPayloadField[any])
}

type customPayloadFieldsAwareInstanaDataObject interface {
	CustomPayloadFieldsAware
	InstanaDataObject
}

// NewCustomPayloadFieldsUnmarshallerAdapter creates a new Unmarshaller instance which can be added as an adapter to the default unmarchallers to map custom payload fields
func NewCustomPayloadFieldsUnmarshallerAdapter[T customPayloadFieldsAwareInstanaDataObject](unmarshaller JSONUnmarshaller[T]) JSONUnmarshaller[T] {
	return &customPayloadFieldsUnmarshallerAdapter[T]{unmarshaller: unmarshaller}
}

type customPayloadFieldsUnmarshallerAdapter[T customPayloadFieldsAwareInstanaDataObject] struct {
	unmarshaller JSONUnmarshaller[T]
}

// UnmarshalArray Unmarshaller interface implementation
func (a *customPayloadFieldsUnmarshallerAdapter[T]) UnmarshalArray(data []byte) (*[]T, error) {
	temp, err := a.unmarshaller.UnmarshalArray(data)
	if err != nil {
		return temp, err
	}
	if temp != nil {
		for _, v := range *temp {
			a.mapCustomPayloadFields(v)
		}
	}
	return temp, nil
}

func (a *customPayloadFieldsUnmarshallerAdapter[T]) Unmarshal(data []byte) (T, error) {
	temp, err := a.unmarshaller.Unmarshal(data)
	if err != nil {
		return temp, err
	}
	a.mapCustomPayloadFields(temp)
	return temp, nil
}

func (a *customPayloadFieldsUnmarshallerAdapter[T]) mapCustomPayloadFields(temp T) {
	customFields := temp.GetCustomerPayloadFields()
	for i, v := range customFields {
		customFields[i] = a.mapCustomPayloadField(v)
	}
	temp.SetCustomerPayloadFields(customFields)
}

func (a *customPayloadFieldsUnmarshallerAdapter[T]) mapCustomPayloadField(field CustomPayloadField[any]) CustomPayloadField[any] {
	if field.Type == DynamicCustomPayloadType {
		data := field.Value.(map[string]interface{})
		var keyPtr *string
		if val, ok := data["key"]; ok && val != nil {
			key := val.(string)
			keyPtr = &key
		}
		field.Value = DynamicCustomPayloadFieldValue{
			TagName: data["tagName"].(string),
			Key:     keyPtr,
		}
	} else {

		value := field.Value.(string)
		field.Value = StaticStringCustomPayloadFieldValue(value)
	}
	return field
}
