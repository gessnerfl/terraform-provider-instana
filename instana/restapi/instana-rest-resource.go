package restapi

//InstanaDataObject is a marker interface for any data object provided by any resource of the Instana REST API
type InstanaDataObject interface {
	GetID() string
	Validate() error
}

//RestResource interface definition of a instana REST resource.
type RestResource interface {
	GetOne(id string) (InstanaDataObject, error)
	Create(data InstanaDataObject) (InstanaDataObject, error)
	Update(data InstanaDataObject) (InstanaDataObject, error)
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

func (r *genericRestResource) Update(data InstanaDataObject) (InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.Put(data, r.resourcePath)
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *genericRestResource) Create(data InstanaDataObject) (InstanaDataObject, error) {
	// Most API endpoints handle a put to new resources a create
	return r.Update(data)
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

// Setup a Rest Resource that uses POST to the collections to create new Resources and
// POST to the individual Item to update. This is necessary since the
// API design for Application Alerts is not consistent with the rest of the Instana API
func NewPostingRestResource(resourcePath string, unmarshaller Unmarshaller, client RestClient) RestResource {
	return &postingRestResource{
		genericRestResource {
			resourcePath: resourcePath,
			unmarshaller: unmarshaller,
			client:       client,
		},
	}
}

type postingRestResource struct {
	genericRestResource
}

func (r *postingRestResource) Update(data InstanaDataObject) (InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.Post(data, r.resourcePath+"/"+data.GetID())
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}

func (r *postingRestResource) Create(data InstanaDataObject) (InstanaDataObject, error) {
	if err := data.Validate(); err != nil {
		return data, err
	}
	response, err := r.client.Post(data, r.resourcePath)
	if err != nil {
		return data, err
	}
	return r.validateResponseAndConvertToStruct(response)
}