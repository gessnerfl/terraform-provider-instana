package restapi

import "encoding/json"

//NewApplicationConfigUnmarshaller creates a new Unmarshaller instance for application configs
func NewApplicationAlertConfigsUnmarshaller() Unmarshaller {
	return &applicationAlertConfigsUnmarshaller{}
}

type applicationAlertConfigsUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *applicationAlertConfigsUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	temp := ApplicationAlertConfigs{

	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return ApplicationAlertConfigs{}, err
	}
	return temp, nil
}

