package jsonconvert

import (
	"errors"
	"fmt"
	"strings"

	"log"
)

const (
	// JSONCONVERTKEY is prefix key
	JSONCONVERTKEY = "$$"
	// JSONSPLITKEY is split key
	JSONSPLITKEY = "."
	// JSONNTIMESKEY is array key
	JSONNTIMESKEY = "[$$n]"
	// RECORDSETKEY is main recotrdset key
	RECORDSETKEY = "$$recordset"
	// JOINRECORDSETKEY is join recotrdset key
	JOINRECORDSETKEY = "$$joinrecordset"
	// JOINRECORDCOLUMNKEY is join column key
	JOINRECORDCOLUMNKEY = "$$joinrecordcolumn"
)

type JsonConverter struct {
	convertList map[string]string
}

func NewJsonConverter() *JsonConverter {
	return &JsonConverter{
		convertList: make(map[string]string, 0),
	}
}

func (jc JsonConverter) SetResponse(identifier, responseJson string) {
	jc.convertList[identifier] = responseJson
}

func (jc JsonConverter) Convert(responsedata []byte, identifier string) ([]byte, error) {
	return jc.convertResponseToJson(jc.convertList[identifier], responsedata)
}

func (jc JsonConverter) convertResponseToJson(convertdata string, responsedata []byte) ([]byte, error) {
	convertdatajson, err := byteToJson([]byte(convertdata))
	if err != nil {
		return nil, err
	}
	jsondata, err := newJSONData(responsedata)
	if err != nil {
		return nil, err
	}
	resultjson, err := jc.convertJsonDataFromResponse(convertdatajson, jsondata)
	if err != nil {
		return responsedata, err
	}
	return jsonToByte(resultjson)
}

func (jc JsonConverter) convertJsonDataFromResponse(convert interface{}, jsondata *jsonData) (interface{}, error) {
	switch convertdata := convert.(type) {
	case map[string]interface{}:
		for key, val := range convertdata {
			if result, err := jc.convertJsonDataFromResponse(val, jsondata); err != nil {
				return nil, err
			} else {
				convertdata[key] = result
			}
		}
		return convertdata, nil
	case []interface{}:
		if jc.isRecordsetConvertData(convertdata) {
			return jc.getRecordsetData(convertdata, jsondata)
		} else {
			resultdata := make([]interface{}, 0)
			for _, dataone := range convertdata {
				result, err := jc.convertJsonDataFromResponse(dataone, jsondata)
				if err != nil {
					return nil, err
				}
				switch data := result.(type) {
				case []interface{}:
					resultdata = append(resultdata, data...)
				default:
					resultdata = append(resultdata, data)
				}
			}
			return resultdata, nil
		}
	case string:
		if strings.HasPrefix(convertdata, JSONCONVERTKEY) == true {
			querykey := strings.Replace(convertdata, JSONCONVERTKEY, "", 1)
			if jc.isNTimeArrayConvertData(querykey) == true {
				return jc.getNTimeArrayData(querykey, jsondata)
			} else {
				if jsonval, err := jsondata.query(querykey); err != nil {
					return "", err
				} else {
					return jsonval, nil
				}
			}
		}
	}
	return convert, nil
}

func (jc JsonConverter) isRecordsetConvertData(convertdata []interface{}) bool {
	for _, convertdataone := range convertdata {
		switch convertdataonedata := convertdataone.(type) {
		case map[string]interface{}:
			if _, ok := convertdataonedata[RECORDSETKEY].(string); ok {
				return true
			}
		}
	}
	return false
}

func (jc JsonConverter) isNTimeArrayConvertData(querykey string) bool {
	if strings.Contains(querykey, JSONNTIMESKEY) == true {
		return true
	}
	return false
}

func (jc JsonConverter) getRecordsetData(convertdata []interface{}, jsondata *jsonData) (interface{}, error) {
	resultdata := make([]interface{}, 0)
	dataset := make([]interface{}, 0)
	joindataset := make([]interface{}, 0)
	var joincolumn string
	for _, convertdataone := range convertdata {
		var err error
		switch convertdataonedata := convertdataone.(type) {
		case map[string]interface{}:
			if dataset, err = jc.getRecordset(RECORDSETKEY, convertdataonedata, jsondata); err != nil {
				return nil, err
			}
			if joindataset, err = jc.getRecordset(JOINRECORDSETKEY, convertdataonedata, jsondata); err != nil {
				log.Print(err)
			}
			if len(joindataset) > 0 {
				var ok bool
				if joincolumn, ok = convertdataonedata[JOINRECORDCOLUMNKEY].(string); ok == false {
					return nil, errors.New("No Join Key")
				}
			}
			delete(convertdataonedata, RECORDSETKEY)
			delete(convertdataonedata, JOINRECORDSETKEY)
			delete(convertdataonedata, JOINRECORDCOLUMNKEY)
			for _, datarecord := range dataset {
				datarecordparser, err := newJSONData(datarecord)
				if err != nil {
					return nil, err
				}
				resultjsonmap := make(map[string]interface{}, 0)
				var joincolumnvalue interface{}
				for key, val := range convertdataonedata {
					switch data := val.(type) {
					case string:
						if data != JOINRECORDSETKEY && strings.HasPrefix(data, JSONCONVERTKEY) == true {
							querykey := strings.Replace(data, JSONCONVERTKEY, "", 1)
							if jsonval, err := datarecordparser.query(querykey); err != nil {
								log.Print(err)
							} else {
								resultjsonmap[key] = jsonval
								if querykey == joincolumn {
									joincolumnvalue = jsonval
								}
							}
						}
					case map[string]interface{}:
						datarecordparser, err := newJSONData(datarecord)
						if err != nil {
							return nil, err
						}
						resultjsonmap[key], err = jc.convertJsonDataFromResponse(data, datarecordparser)
						if err != nil {
							return nil, err
						}
					}
				}
				if len(joindataset) > 0 {
					for key, val := range convertdataonedata {
						switch data := val.(type) {
						case string:
							if data == JOINRECORDSETKEY {
								resultjsonmap[key] = jc.getDataFromJoinRecordset(joincolumn, joincolumnvalue, joindataset)
							}
						}
					}
				}
				resultdata = append(resultdata, resultjsonmap)
			}
		}
	}
	return resultdata, nil
}

func (jc JsonConverter) getNTimeArrayData(querykey string, jsondata *jsonData) (interface{}, error) {
	splitquerylist := make([]string, 0)
	log.Print(querykey)
	if strings.HasPrefix(querykey, JSONNTIMESKEY) == true {
		splitquerylist = strings.Split(querykey, fmt.Sprintf("%s%s", JSONNTIMESKEY, JSONSPLITKEY))
		splitquerylist[0] = JSONSPLITKEY
	} else {
		splitquerylist = strings.Split(querykey, fmt.Sprintf("%s%s%s", JSONSPLITKEY, JSONNTIMESKEY, JSONSPLITKEY))
	}
	log.Print(splitquerylist)
	if jsonarray, err := jsondata.query(splitquerylist[0]); err != nil {
		log.Print(err)
		return nil, err
	} else {
		datalist := make([]interface{}, 0)
		switch jsonarraydata := jsonarray.(type) {
		case map[string]interface{}:
			if valone, ok := jsonarraydata[splitquerylist[1]]; ok {
				return valone, nil
			}
		case []interface{}:
			for _, jsonarrayval := range jsonarraydata {
				arrayjsonparser, err := newJSONData(jsonarrayval)
				if err != nil {
					return nil, err
				}
				if jsonval, err := arrayjsonparser.query(splitquerylist[1]); err != nil {
					return nil, err
				} else {
					datalist = append(datalist, jsonval)
				}
			}
			return datalist, nil
		}
	}
	return nil, errors.New("No Convert")
}

func (jc JsonConverter) getRecordset(querykey string, convertdata map[string]interface{}, jsonparser *jsonData) ([]interface{}, error) {
	if datasetquerykey, ok := convertdata[querykey].(string); ok {
		if jsonval, err := jsonparser.query(datasetquerykey); err != nil {
			log.Print(err)
			return nil, err
		} else {
			switch jsonvaldata := jsonval.(type) {
			case []interface{}:
				return jsonvaldata, nil
			case interface{}:
				return []interface{}{jsonvaldata}, nil
			}
		}
	}
	return nil, fmt.Errorf("No RecordSet %s", querykey)
}

func (jc JsonConverter) getDataFromJoinRecordset(joincolumn string, joincolumnvalue interface{}, joindata []interface{}) interface{} {
	for _, joindataone := range joindata {
		switch joindataonedata := joindataone.(type) {
		case map[string]interface{}:
			if columdata, ok := joindataonedata[joincolumn]; ok {
				if columdata == joincolumnvalue {
					return joindataone
				}
			}
		}
	}
	return nil
}
