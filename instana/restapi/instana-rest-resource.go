package restapi

//InstanaDataObject is a marker interface for any data object provided by any resource of the Instana REST API
type InstanaDataObject interface {
	GetIDForResourcePath() string
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

//DataFilterFunc function definition for filtering data received from Instana API
type DataFilterFunc func(o InstanaDataObject) bool

//ReadOnlyRestResource interface definition for a read only REST resource. The resource at instana might
//implement more methods but the implementation of the provider is limited to read only.
type ReadOnlyRestResource interface {
	GetAll() (*[]InstanaDataObject, error)
	GetOne(id string) (InstanaDataObject, error)
}

//JSONUnmarshaller interface definition for unmarshalling that unmarshalls JSON to go data structures
type JSONUnmarshaller interface {
	//Unmarshal converts the provided json bytes into the go data data structure as provided in the target
	Unmarshal(data []byte) (interface{}, error)
}
