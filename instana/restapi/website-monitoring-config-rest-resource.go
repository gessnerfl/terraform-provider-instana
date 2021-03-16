package restapi

import "errors"

//NewWebsiteMonitoringConfigRestResource creates a new REST for the website monitoring config
func NewWebsiteMonitoringConfigRestResource(unmarshaller JSONUnmarshaller, client RestClient) RestResource {
	return &websiteMonitoringConfigRestResource{
		resourcePath: WebsiteMonitoringConfigResourcePath,
		unmarshaller: unmarshaller,
		client:       client,
	}
}

type websiteMonitoringConfigRestResource struct {
	resourcePath string
	unmarshaller JSONUnmarshaller
	client       RestClient
}

func (r *websiteMonitoringConfigRestResource) GetOne(id string) (InstanaDataObject, error) {
	data, err := r.client.GetOne(id, r.resourcePath)
	if err != nil {
		return nil, err
	}
	return r.validateResponseAndConvertToStruct(data)
}

func (r *websiteMonitoringConfigRestResource) Create(data InstanaDataObject) (InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.PostByQuery(r.resourcePath, map[string]string{"name": data.(*WebsiteMonitoringConfig).Name})
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *websiteMonitoringConfigRestResource) Update(data InstanaDataObject) (InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.PutByQuery(r.resourcePath, data.GetIDForResourcePath(), map[string]string{"name": data.(*WebsiteMonitoringConfig).Name})
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *websiteMonitoringConfigRestResource) validateResponseAndConvertToStruct(data []byte) (InstanaDataObject, error) {
	object, err := r.unmarshaller.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	dataObject, ok := object.(InstanaDataObject)
	if !ok {
		return dataObject, errors.New("Unmarshalled object does not implement InstanaDataObject")
	}

	if err := dataObject.Validate(); err != nil {
		return dataObject, err
	}
	return dataObject, nil
}

func (r *websiteMonitoringConfigRestResource) Delete(data InstanaDataObject) error {
	return r.DeleteByID(data.GetIDForResourcePath())
}

func (r *websiteMonitoringConfigRestResource) DeleteByID(id string) error {
	return r.client.Delete(id, r.resourcePath)
}
