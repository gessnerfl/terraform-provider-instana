package restapi

// NewSyntheticTestRestResource creates a new REST resource using the provided unmarshaller function to convert the response from the REST API to the corresponding InstanaDataObject. The REST resource is using PUT as operation for create and update
func NewSyntheticTestRestResource(unmarshaller JSONUnmarshaller[*SyntheticTest], client RestClient) RestResource[*SyntheticTest] {
	return &SyntheticTestRestResource{
		resourcePath: SyntheticTestResourcePath,
		unmarshaller: unmarshaller,
		client:       client,
	}
}

type SyntheticTestRestResource struct {
	resourcePath string
	unmarshaller JSONUnmarshaller[*SyntheticTest]
	client       RestClient
}

func (r *SyntheticTestRestResource) GetAll() (*[]*SyntheticTest, error) {
	data, err := r.client.Get(r.resourcePath)
	if err != nil {
		return nil, err
	}
	objects, err := r.unmarshaller.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}
	return objects, nil
}

func (r *SyntheticTestRestResource) GetOne(id string) (*SyntheticTest, error) {
	data, err := r.client.GetOne(id, r.resourcePath)
	if err != nil {
		return nil, err
	}
	return r.validateResponseAndConvertToStruct(data)
}

func (r *SyntheticTestRestResource) Create(data *SyntheticTest) (*SyntheticTest, error) {
	response, err := r.client.Post(data, r.resourcePath)
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *SyntheticTestRestResource) Update(data *SyntheticTest) (*SyntheticTest, error) {
	_, err := r.client.Put(data, r.resourcePath)
	if err != nil {
		return data, err
	}
	return r.GetOne(data.GetIDForResourcePath())
}

func (r *SyntheticTestRestResource) validateResponseAndConvertToStruct(data []byte) (*SyntheticTest, error) {
	dataObject, err := r.unmarshaller.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return dataObject, nil
}

func (r *SyntheticTestRestResource) Delete(data *SyntheticTest) error {
	return r.DeleteByID(data.GetIDForResourcePath())
}

func (r *SyntheticTestRestResource) DeleteByID(id string) error {
	return r.client.Delete(id, r.resourcePath)
}
