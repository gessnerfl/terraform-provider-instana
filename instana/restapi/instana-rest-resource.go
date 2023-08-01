package restapi

// InstanaDataObject is a marker interface for any data object provided by any resource of the Instana REST API
type InstanaDataObject interface {
	GetIDForResourcePath() string
	Validate() error
}

// RestResource interface definition of a instana REST resource.
type RestResource[T InstanaDataObject] interface {
	ReadOnlyRestResource[T]
	Create(data T) (T, error)
	Update(data T) (T, error)
	Delete(data T) error
	DeleteByID(id string) error
}

// DataFilterFunc function definition for filtering data received from Instana API
type DataFilterFunc func(o InstanaDataObject) bool

// ReadOnlyRestResource interface definition for a read only REST resource. The resource at instana might
// implement more methods but the implementation of the provider is limited to read only.
type ReadOnlyRestResource[T InstanaDataObject] interface {
	GetAll() (*[]T, error)
	GetOne(id string) (T, error)
}

// JSONUnmarshaller interface definition for unmarshalling that unmarshalls JSON to go data structures
type JSONUnmarshaller[T any] interface {
	//Unmarshal converts the provided json bytes into the go data structure as provided in the target
	Unmarshal(data []byte) (T, error)
	//UnmarshalArray converts the provided json bytes into the go data structure as provided in the target
	UnmarshalArray(data []byte) (*[]T, error)
}
