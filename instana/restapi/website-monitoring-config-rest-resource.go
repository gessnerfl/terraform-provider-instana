package restapi

//NewWebsiteMonitoringConfigRestResource creates a new REST for the website monitoring config
func NewWebsiteMonitoringConfigRestResource(unmarshaller Unmarshaller, client RestClient) RestResource {
	return &websiteMonitoringConfigRestResource{
		resourcePath: WebsiteMonitoringConfigResourcePath,
		unmarshaller: unmarshaller,
		client:       client,
	}
}

type websiteMonitoringConfigRestResource struct {
	resourcePath string
	unmarshaller Unmarshaller
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
	response, err := r.client.PostByQuery(r.resourcePath, map[string]string{"name": data.(WebsiteMonitoringConfig).Name})
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *websiteMonitoringConfigRestResource) Update(data InstanaDataObject) (InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.PutByQuery(r.resourcePath, data.GetID(), map[string]string{"name": data.(WebsiteMonitoringConfig).Name})
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *websiteMonitoringConfigRestResource) validateResponseAndConvertToStruct(data []byte) (InstanaDataObject, error) {
	object, err := r.unmarshaller.Unmarshal(data)
	if err != nil {
		return object, err
	}

	if err := object.Validate(); err != nil {
		return object, err
	}
	return object, nil
}

func (r *websiteMonitoringConfigRestResource) Delete(data InstanaDataObject) error {
	return r.DeleteByID(data.GetID())
}

func (r *websiteMonitoringConfigRestResource) DeleteByID(id string) error {
	return r.client.Delete(id, r.resourcePath)
}
