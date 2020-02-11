package restapi

//InstanaDataObject is a marker interface for any data object provided by any resource of the Instana REST API
type InstanaDataObject interface {
	GetID() string
	Validate() error
}

//RestResource interface definition of a instana REST resource.
type RestResource interface {
	GetOne(id string) (InstanaDataObject, error)
	Upsert(data InstanaDataObject) (InstanaDataObject, error)
	Delete(data InstanaDataObject) error
	DeleteByID(id string) error
}

//Unmarshaller interface definition for unmarshalling the binary data to the desired struct
type Unmarshaller interface {
	Unmarshal(data []byte) (InstanaDataObject, error)
}

//NewRestResource creates a new REST resource using the provided unmarshaller function to convert the response from the REST API to the corresponding InstanaDataObject
func NewRestResource(resourcePath string, unmarshaller Unmarshaller, client RestClient) RestResource {
	return &genericRestResource{
		resourcePath: resourcePath,
		unmarshaller: unmarshaller,
		client:       client,
	}
}

type genericRestResource struct {
	resourcePath string
	unmarshaller Unmarshaller
	client       RestClient
}

func (r *genericRestResource) GetOne(id string) (InstanaDataObject, error) {
	data, err := r.client.GetOne(id, r.resourcePath)
	if err != nil {
		return nil, err
	}
	return r.validateResponseAndConvertToStruct(data)
}

func (r *genericRestResource) Upsert(data InstanaDataObject) (InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.Put(data, r.resourcePath)
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *genericRestResource) validateResponseAndConvertToStruct(data []byte) (InstanaDataObject, error) {
	object, err := r.unmarshaller.Unmarshal(data)
	if err != nil {
		return object, err
	}

	if err := object.Validate(); err != nil {
		return object, err
	}
	return object, nil
}

func (r *genericRestResource) Delete(data InstanaDataObject) error {
	return r.DeleteByID(data.GetID())
}

func (r *genericRestResource) DeleteByID(id string) error {
	return r.client.Delete(id, r.resourcePath)
}
