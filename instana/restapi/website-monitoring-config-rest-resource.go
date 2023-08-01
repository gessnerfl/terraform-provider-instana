package restapi

// NewWebsiteMonitoringConfigRestResource creates a new REST for the website monitoring config
func NewWebsiteMonitoringConfigRestResource(unmarshaller JSONUnmarshaller[*WebsiteMonitoringConfig], client RestClient) RestResource[*WebsiteMonitoringConfig] {
	return &websiteMonitoringConfigRestResource{
		resourcePath: WebsiteMonitoringConfigResourcePath,
		unmarshaller: unmarshaller,
		client:       client,
	}
}

type websiteMonitoringConfigRestResource struct {
	resourcePath string
	unmarshaller JSONUnmarshaller[*WebsiteMonitoringConfig]
	client       RestClient
}

func (r *websiteMonitoringConfigRestResource) GetOne(id string) (*WebsiteMonitoringConfig, error) {
	data, err := r.client.GetOne(id, r.resourcePath)
	if err != nil {
		return nil, err
	}
	return r.validateResponseAndConvertToStruct(data)
}

func (r *websiteMonitoringConfigRestResource) Create(data *WebsiteMonitoringConfig) (*WebsiteMonitoringConfig, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.PostByQuery(r.resourcePath, map[string]string{"name": data.Name})
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *websiteMonitoringConfigRestResource) Update(data *WebsiteMonitoringConfig) (*WebsiteMonitoringConfig, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.PutByQuery(r.resourcePath, data.GetIDForResourcePath(), map[string]string{"name": data.Name})
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *websiteMonitoringConfigRestResource) validateResponseAndConvertToStruct(data []byte) (*WebsiteMonitoringConfig, error) {
	dataObject, err := r.unmarshaller.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	if err := dataObject.Validate(); err != nil {
		return dataObject, err
	}
	return dataObject, nil
}

func (r *websiteMonitoringConfigRestResource) Delete(data *WebsiteMonitoringConfig) error {
	return r.DeleteByID(data.GetIDForResourcePath())
}

func (r *websiteMonitoringConfigRestResource) DeleteByID(id string) error {
	return r.client.Delete(id, r.resourcePath)
}
