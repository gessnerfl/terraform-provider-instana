package restapi

import "encoding/json"

func unmarshalArray[T InstanaDataObject](data []byte, unmarshalObject func([]byte) (T, error)) (*[]T, error) {
	rawArray := make([]json.RawMessage, 0)
	if err := json.Unmarshal(data, &rawArray); err != nil {
		return nil, err
	}

	result := make([]T, len(rawArray))
	for i, r := range rawArray {
		c, err := unmarshalObject(r)
		if err != nil {
			return nil, err
		}
		result[i] = c
	}
	return &result, nil
}
