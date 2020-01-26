package services

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

type unmarshallingFunc func(data []byte) (restapi.InstanaDataObject, error)

//NewRestResource creates a new REST resource using the provided unmarshalling function to convert the response from the REST API to the corresponding InstanaDataObject
func NewRestResource(resourcePath string, unmarshallingFunc unmarshallingFunc, client restapi.RestClient) restapi.RestResource {
	return &genericRestResource{
		resourcePath:      resourcePath,
		unmarshallingFunc: unmarshallingFunc,
		client:            client,
	}
}

type genericRestResource struct {
	resourcePath      string
	unmarshallingFunc unmarshallingFunc
	client            restapi.RestClient
}

func (r *genericRestResource) GetOne(id string) (restapi.InstanaDataObject, error) {
	data, err := r.client.GetOne(id, r.resourcePath)
	if err != nil {
		return nil, err
	}
	return r.validateResponseAndConvertToStruct(data)
}

func (r *genericRestResource) Upsert(data restapi.InstanaDataObject) (restapi.InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.Put(data, r.resourcePath)
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *genericRestResource) validateResponseAndConvertToStruct(data []byte) (restapi.InstanaDataObject, error) {
	object, err := r.unmarshallingFunc(data)
	if err != nil {
		return object, err
	}

	if err := object.Validate(); err != nil {
		return object, err
	}
	return object, nil
}

func (r *genericRestResource) Delete(data restapi.InstanaDataObject) error {
	return r.DeleteByID(data.GetID())
}

func (r *genericRestResource) DeleteByID(id string) error {
	return r.client.Delete(id, r.resourcePath)
}
