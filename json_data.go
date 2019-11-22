package jsonconvert

import (
	"fmt"
	"strconv"
	"strings"
)

type jsonData struct {
	Data interface{}
}

func newJSONData(inputdata interface{}) (*jsonData, error) {
	var data = new(interface{})
	switch converteddata := inputdata.(type) {
	case []byte:
		err := byteToJsonStruct(converteddata, data)
		if err != nil {
			return nil, err
		}
		return &jsonData{*data}, nil
	case string:
		err := byteToJsonStruct([]byte(converteddata), data)
		if err != nil {
			return nil, err
		}
		return &jsonData{*data}, nil
	default:
		return &jsonData{Data: inputdata}, nil
	}
}

func (jd *jsonData) query(exp string) (interface{}, error) {
	if exp == "." {
		return jd.Data, nil
	}
	paths := strings.Split(exp, ".")
	var context interface{} = jd.Data
	for _, path := range paths {
		if len(path) >= 3 && strings.HasPrefix(path, "[") && strings.HasSuffix(path, "]") {
			// array
			index, err := strconv.Atoi(path[1 : len(path)-1])
			if err != nil {
				return nil, err
			}
			if v, ok := context.([]interface{}); ok {
				if len(v) <= index {
					return nil, fmt.Errorf("%s: index out of range", path)
				}
				context = v[index]
			} else {
				return nil, fmt.Errorf("%s is not an array. %v", path, v)
			}
		} else {
			// map
			if v, ok := context.(map[string]interface{}); ok {
				if val, ok := v[path]; ok {
					context = val
				} else {
					return nil, fmt.Errorf("%s does not exist", path)
				}
			} else {
				return nil, fmt.Errorf("%s is not an object. %v", path, v)
			}
		}
	}
	switch converteddata := context.(type) {
	case map[string]interface{}:
		return converteddata, nil
	case []interface{}:
		return converteddata, nil
	case string:
		return converteddata, nil
	case int:
		return converteddata, nil
	case int32:
		return converteddata, nil
	case int64:
		return converteddata, nil
	case float32:
		return converteddata, nil
	case float64:
		return converteddata, nil
	case bool:
		return converteddata, nil
	default:
		return converteddata, nil
	}
}
