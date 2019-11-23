package jsondata

import (
	"fmt"
	"github.com/howood/jsonconvert/internal/parser"
	"strconv"
	"strings"
)

// JSONData represents JSON data entity
type JSONData struct {
	// Data is JSON Data
	Data interface{}
}

// NewJSONData create JSONData pointer from []byte / string / interface{}
func NewJSONData(inputdata interface{}) (*JSONData, error) {
	var data = new(interface{})
	switch converteddata := inputdata.(type) {
	case []byte:
		err := parser.ByteToJsonStruct(converteddata, data)
		if err != nil {
			return nil, err
		}
		return &JSONData{*data}, nil
	case string:
		err := parser.ByteToJsonStruct([]byte(converteddata), data)
		if err != nil {
			return nil, err
		}
		return &JSONData{*data}, nil
	default:
		return &JSONData{Data: inputdata}, nil
	}
}

// Query is extract data from JSON with item key
func (jd *JSONData) Query(key string) (interface{}, error) {
	if key == "." {
		return jd.Data, nil
	}
	paths := strings.Split(key, ".")
	var context interface{} = jd.Data
	for _, path := range paths {
		if len(path) >= 3 && strings.HasPrefix(path, "[") && strings.HasSuffix(path, "]") {
			// array
			index, err := strconv.Atoi(path[1 : len(path)-1])
			if err != nil {
				return nil, err
			}
			switch v := context.(type) {
			case []interface{}:
				if len(v) <= index {
					return nil, fmt.Errorf("%s: index out of range", path)
				}
				context = v[index]
			default:
				return nil, fmt.Errorf("%s: not array. %v", path, v)
			}
		} else {
			// map
			switch v := context.(type) {
			case map[string]interface{}:
				if val, ok := v[path]; ok {
					context = val
				} else {
					return nil, fmt.Errorf("%s: not exist", path)
				}
			default:
				return nil, fmt.Errorf("%s: not object. %v", path, v)
			}
		}
	}
	return context, nil
}
