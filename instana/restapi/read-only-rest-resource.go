package restapi

import (
	"reflect"
)

//NewReadOnlyRestResource creates a new instance of ReadOnlyRestResource
func NewReadOnlyRestResource(resourcePath string, objectUnmarshaller JSONUnmarshaller, arrayUnmarshaller JSONUnmarshaller, client RestClient) ReadOnlyRestResource {
	return &readOnlyRestResource{
		resourcePath:       resourcePath,
		objectUnmarshaller: objectUnmarshaller,
		arrayUnmarshaller:  arrayUnmarshaller,
		client:             client,
	}
}

type readOnlyRestResource struct {
	resourcePath       string
	objectUnmarshaller JSONUnmarshaller
	arrayUnmarshaller  JSONUnmarshaller
	client             RestClient
}

func (i *readOnlyRestResource) GetAll() (*[]InstanaDataObject, error) {
	data, err := i.client.Get(i.resourcePath)
	if err != nil {
		return nil, err
	}
	objects, err := i.arrayUnmarshaller.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	values := reflect.ValueOf(objects).Elem()
	result := make([]InstanaDataObject, values.Len())

	for i := 0; i < values.Len(); i++ {
		o := values.Index(i).Interface()
		result[i] = o.(InstanaDataObject)
	}
	return &result, nil
}

func (i *readOnlyRestResource) GetOne(id string) (InstanaDataObject, error) {
	data, err := i.client.GetOne(id, i.resourcePath)
	if err != nil {
		return nil, err
	}
	object, err := i.objectUnmarshaller.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(object).Elem()
	return value.Interface().(InstanaDataObject), nil
}
