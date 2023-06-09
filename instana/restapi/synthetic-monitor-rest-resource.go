package restapi

import "errors"

// NewSyntheticMonitorRestResource creates a new REST resource using the provided unmarshaller function to convert the response from the REST API to the corresponding InstanaDataObject. The REST resource is using PUT as operation for create and update
func NewSyntheticMonitorRestResource(unmarshaller JSONUnmarshaller, client RestClient) RestResource {
	return &syntheticMonitorRestResource{
		resourcePath: SyntheticMonitorResourcePath,
		unmarshaller: unmarshaller,
		client:       client,
	}
}

type syntheticMonitorRestResource struct {
	resourcePath string
	unmarshaller JSONUnmarshaller
	client       RestClient
}

func (r *syntheticMonitorRestResource) GetOne(id string) (InstanaDataObject, error) {
	data, err := r.client.GetOne(id, r.resourcePath)
	if err != nil {
		return nil, err
	}
	return r.validateResponseAndConvertToStruct(data)
}

func (r *syntheticMonitorRestResource) Create(data InstanaDataObject) (InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.Post(data, r.resourcePath)
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *syntheticMonitorRestResource) Update(data InstanaDataObject) (InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	_, err := r.client.Put(data, r.resourcePath)
	if err != nil {
		return data, err
	}
	return r.GetOne(data.GetIDForResourcePath())
}

func (r *syntheticMonitorRestResource) validateResponseAndConvertToStruct(data []byte) (InstanaDataObject, error) {
	object, err := r.unmarshaller.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	dataObject, ok := object.(InstanaDataObject)
	if !ok {
		return dataObject, errors.New("unmarshalled object does not implement InstanaDataObject")
	}

	if err := dataObject.Validate(); err != nil {
		return dataObject, err
	}
	return dataObject, nil
}

func (r *syntheticMonitorRestResource) Delete(data InstanaDataObject) error {
	return r.DeleteByID(data.GetIDForResourcePath())
}

func (r *syntheticMonitorRestResource) DeleteByID(id string) error {
	return r.client.Delete(id, r.resourcePath)
}
