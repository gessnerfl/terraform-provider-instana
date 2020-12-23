package restapi

//NewPUTOnlyRestResource creates a new REST resource using the provided unmarshaller function to convert the response from the REST API to the corresponding InstanaDataObject. The REST resource is using PUT as operation for create and update
func NewPUTOnlyRestResource(resourcePath string, unmarshaller Unmarshaller, client RestClient) RestResource {
	return &putOnlyRestResource{
		resourcePath: resourcePath,
		unmarshaller: unmarshaller,
		client:       client,
	}
}

type putOnlyRestResource struct {
	resourcePath string
	unmarshaller Unmarshaller
	client       RestClient
}

func (r *putOnlyRestResource) GetOne(id string) (InstanaDataObject, error) {
	data, err := r.client.GetOne(id, r.resourcePath)
	if err != nil {
		return nil, err
	}
	return r.validateResponseAndConvertToStruct(data)
}

func (r *putOnlyRestResource) Create(data InstanaDataObject) (InstanaDataObject, error) {
	return r.Update(data)
}

func (r *putOnlyRestResource) Update(data InstanaDataObject) (InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.Put(data, r.resourcePath)
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *putOnlyRestResource) validateResponseAndConvertToStruct(data []byte) (InstanaDataObject, error) {
	object, err := r.unmarshaller.Unmarshal(data)
	if err != nil {
		return object, err
	}

	if err := object.Validate(); err != nil {
		return object, err
	}
	return object, nil
}

func (r *putOnlyRestResource) Delete(data InstanaDataObject) error {
	return r.DeleteByID(data.GetID())
}

func (r *putOnlyRestResource) DeleteByID(id string) error {
	return r.client.Delete(id, r.resourcePath)
}
