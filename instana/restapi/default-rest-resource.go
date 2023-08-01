package restapi

import (
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

// NewCreatePUTUpdatePUTRestResource creates a new REST resource using the provided unmarshaller function to convert the response from the REST API to the corresponding InstanaDataObject. The REST resource is using PUT as operation for create and update
func NewCreatePUTUpdatePUTRestResource[T InstanaDataObject](resourcePath string, unmarshaller JSONUnmarshaller[T], client RestClient) RestResource[T] {
	return &defaultRestResource[T]{
		mode:         DefaultRestResourceModeCreateAndUpdatePUT,
		resourcePath: resourcePath,
		unmarshaller: unmarshaller,
		client:       client,
	}
}

// NewCreatePOSTUpdatePUTRestResource creates a new REST resource using the provided unmarshaller function to convert the response from the REST API to the corresponding InstanaDataObject. The REST resource is using POST as operation for create and PUT for update
func NewCreatePOSTUpdatePUTRestResource[T InstanaDataObject](resourcePath string, unmarshaller JSONUnmarshaller[T], client RestClient) RestResource[T] {
	return &defaultRestResource[T]{
		mode:         DefaultRestResourceModeCreatePOSTUpdatePUT,
		resourcePath: resourcePath,
		unmarshaller: unmarshaller,
		client:       client,
	}
}

// NewCreatePOSTUpdatePOSTRestResource creates a new REST resource using the provided unmarshaller function to convert the response from the REST API to the corresponding InstanaDataObject. The REST resource is using POST as operation for create and update
func NewCreatePOSTUpdatePOSTRestResource[T InstanaDataObject](resourcePath string, unmarshaller JSONUnmarshaller[T], client RestClient) RestResource[T] {
	return &defaultRestResource[T]{
		mode:         DefaultRestResourceModeCreateAndUpdatePOST,
		resourcePath: resourcePath,
		unmarshaller: unmarshaller,
		client:       client,
	}
}

// DefaultRestResourceMode custom type for create/update behavior of the defaultRestResource
type DefaultRestResourceMode string

type restClientOperation func(InstanaDataObject, string) ([]byte, error)

const (
	//DefaultRestResourceModeCreateAndUpdatePUT constant value for the DefaultRestResourceMode CREATE_PUT_UPDATE_PUT where create and update is implemented as an upsert using HTTP PUT method only
	DefaultRestResourceModeCreateAndUpdatePUT = DefaultRestResourceMode("CREATE_PUT_UPDATE_PUT")
	//DefaultRestResourceModeCreatePOSTUpdatePUT constant value for the DefaultRestResourceMode CREATE_POST_UPDATE_PUT where create is implemented as an HTTP POST method and update is implemented as HTTP PUT method
	DefaultRestResourceModeCreatePOSTUpdatePUT = DefaultRestResourceMode("CREATE_POST_UPDATE_PUT")
	//DefaultRestResourceModeCreateAndUpdatePOST constant value for the DefaultRestResourceMode CREATE_POST_UPDATE_POST where create is implemented as an HTTP POST method and update is implemented as HTTP PUT method
	DefaultRestResourceModeCreateAndUpdatePOST = DefaultRestResourceMode("CREATE_POST_UPDATE_POST")
)

type defaultRestResource[T InstanaDataObject] struct {
	mode         DefaultRestResourceMode
	resourcePath string
	unmarshaller JSONUnmarshaller[T]
	client       RestClient
}

func (r *defaultRestResource[T]) GetOne(id string) (T, error) {
	data, err := r.client.GetOne(id, r.resourcePath)
	if err != nil {
		return utils.GetZeroValue[T](), err
	}
	return r.validateResponseAndConvertToStruct(data)
}

func (r *defaultRestResource[T]) Create(data T) (T, error) {
	if r.mode == DefaultRestResourceModeCreateAndUpdatePUT {
		return r.upsert(data, r.client.Put)
	}
	return r.upsert(data, r.client.Post)
}

func (r *defaultRestResource[T]) Update(data T) (T, error) {
	if r.mode == DefaultRestResourceModeCreateAndUpdatePOST {
		return r.upsert(data, r.client.PostWithID)
	}
	return r.upsert(data, r.client.Put)
}

func (r *defaultRestResource[T]) upsert(data T, operation restClientOperation) (T, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := operation(data, r.resourcePath)
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *defaultRestResource[T]) validateResponseAndConvertToStruct(data []byte) (T, error) {
	dataObject, err := r.unmarshaller.Unmarshal(data)
	if err != nil {
		return utils.GetZeroValue[T](), err
	}

	if err := dataObject.Validate(); err != nil {
		return dataObject, err
	}
	return dataObject, nil
}

func (r *defaultRestResource[T]) Delete(data T) error {
	return r.DeleteByID(data.GetIDForResourcePath())
}

func (r *defaultRestResource[T]) DeleteByID(id string) error {
	return r.client.Delete(id, r.resourcePath)
}
