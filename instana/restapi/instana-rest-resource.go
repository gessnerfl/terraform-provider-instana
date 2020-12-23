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
