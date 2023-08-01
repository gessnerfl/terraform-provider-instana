package restapi

import "github.com/gessnerfl/terraform-provider-instana/utils"

// NewReadOnlyRestResource creates a new instance of ReadOnlyRestResource
func NewReadOnlyRestResource[T InstanaDataObject](resourcePath string, objectUnmarshaller JSONUnmarshaller[T], arrayUnmarshaller JSONUnmarshaller[*[]T], client RestClient) ReadOnlyRestResource[T] {
	return &readOnlyRestResource[T]{
		resourcePath:       resourcePath,
		objectUnmarshaller: objectUnmarshaller,
		arrayUnmarshaller:  arrayUnmarshaller,
		client:             client,
	}
}

type readOnlyRestResource[T InstanaDataObject] struct {
	resourcePath       string
	objectUnmarshaller JSONUnmarshaller[T]
	arrayUnmarshaller  JSONUnmarshaller[*[]T]
	client             RestClient
}

func (r *readOnlyRestResource[T]) GetAll() (*[]T, error) {
	data, err := r.client.Get(r.resourcePath)
	if err != nil {
		return nil, err
	}
	objects, err := r.arrayUnmarshaller.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return objects, nil
}

func (r *readOnlyRestResource[T]) GetOne(id string) (T, error) {
	data, err := r.client.GetOne(id, r.resourcePath)
	if err != nil {
		return utils.GetZeroValue[T](), err
	}
	return r.objectUnmarshaller.Unmarshal(data)
}
