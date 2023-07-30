package restapi

import (
	"reflect"
)

// NewReadOnlyRestResource creates a new instance of ReadOnlyRestResource
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

func (r *readOnlyRestResource) GetAll() (*[]InstanaDataObject, error) {
	data, err := r.client.Get(r.resourcePath)
	if err != nil {
		return nil, err
	}
	objects, err := r.arrayUnmarshaller.Unmarshal(data)
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

func (r *readOnlyRestResource) GetOne(id string) (InstanaDataObject, error) {
	data, err := r.client.GetOne(id, r.resourcePath)
	if err != nil {
		return nil, err
	}
	object, err := r.objectUnmarshaller.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(object).Elem()
	return value.Interface().(InstanaDataObject), nil
}
