package jsonconvert

import (
	"errors"
	"fmt"
	"strings"

	iJd "github.com/howood/jsonconvert/internal/jsondata"
	iP "github.com/howood/jsonconvert/internal/parser"
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
	// JOINRECORDCOLUMNKEY is join column link key
	JOINRECORDCOLUMNKEY = "$$joinrecordcolumn"
)

// JSONConvert represents JSON convert entity
type JSONConvert struct {
	// convertList is convert list settings
	convertList map[string]string
}

// NewJSONConvert create JSONConvert pointer
func NewJSONConvert() *JSONConvert {
	return &JSONConvert{
		convertList: make(map[string]string, 0),
	}
}

// SetResponseSetting is set convert list response settings
func (jc JSONConvert) SetResponseSetting(identifier, responseJSONSetting string) {
	jc.convertList[identifier] = responseJSONSetting
}

// Convert is convert JSON data to other format
func (jc JSONConvert) Convert(inputdata []byte, identifier string) ([]byte, error) {
	return jc.convertData(jc.convertList[identifier], inputdata)
}

// convertData is convert input data
func (jc JSONConvert) convertData(convertdata string, inputdata []byte) ([]byte, error) {
	convertdatajson, err := iP.ByteToJson([]byte(convertdata))
	if err != nil {
		return nil, err
	}
	jsondata, err := iJd.NewJSONData(inputdata)
	if err != nil {
		return nil, err
	}
	resultjson, err := jc.convertJSONData(convertdatajson, jsondata)
	if err != nil {
		return inputdata, err
	}
	return iP.JsonToByte(resultjson)
}

// convertJSONData is convert using jsondata.JSONData
func (jc JSONConvert) convertJSONData(convert interface{}, jsondata *iJd.JSONData) (interface{}, error) {
	switch convertdata := convert.(type) {
	case map[string]interface{}:
		for key, val := range convertdata {
			result, err := jc.convertJSONData(val, jsondata)
			if err != nil {
				return nil, err
			}
			convertdata[key] = result
		}
		return convertdata, nil
	case []interface{}:
		if jc.isRecordset(convertdata) {
			return jc.getRecordsetData(convertdata, jsondata)
		}
		resultdata := make([]interface{}, 0)
		for _, dataone := range convertdata {
			result, err := jc.convertJSONData(dataone, jsondata)
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
	case string:
		if strings.HasPrefix(convertdata, JSONCONVERTKEY) == true {
			querykey := strings.Replace(convertdata, JSONCONVERTKEY, "", 1)
			if jc.isNTimeArray(querykey) == true {
				return jc.getNTimeArrayData(querykey, jsondata)
			}
			jsonval, err := jsondata.Query(querykey)
			if err != nil {
				return "", err
			}
			return jsonval, nil
		}
	}
	return convert, nil
}

// isRecordset is check Recordset data or not
func (jc JSONConvert) isRecordset(convertdata []interface{}) bool {
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

// isNTimeArray is check NTimeArray data or not
func (jc JSONConvert) isNTimeArray(querykey string) bool {
	if strings.Contains(querykey, JSONNTIMESKEY) == true {
		return true
	}
	return false
}

// getRecordsetData get Recordset data
func (jc JSONConvert) getRecordsetData(convertdata []interface{}, jsondata *iJd.JSONData) (interface{}, error) {
	resultdata := make([]interface{}, 0)
	dataset := make([]interface{}, 0)
	joindataset := make([]interface{}, 0)
	var joincolumn string
	for _, convertdataone := range convertdata {
		var err error
		switch convertdataonedata := convertdataone.(type) {
		case map[string]interface{}:
			if dataset, err = jc.getRecordsetWithQueryKey(RECORDSETKEY, convertdataonedata, jsondata); err != nil {
				return nil, err
			}
			if joindataset, err = jc.getRecordsetWithQueryKey(JOINRECORDSETKEY, convertdataonedata, jsondata); err != nil {
				//log.Print(err)
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
				datarecordparser, err := iJd.NewJSONData(datarecord)
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
							jsonval, err := datarecordparser.Query(querykey)
							if err != nil {
								//log.Print(err)
							}
							resultjsonmap[key] = jsonval
							if querykey == joincolumn {
								joincolumnvalue = jsonval
							}
						}
					case map[string]interface{}:
						datarecordparser, err := iJd.NewJSONData(datarecord)
						if err != nil {
							return nil, err
						}
						resultjsonmap[key], err = jc.convertJSONData(data, datarecordparser)
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

// getNTimeArrayData get NTimeArray data
func (jc JSONConvert) getNTimeArrayData(querykey string, jsondata *iJd.JSONData) (interface{}, error) {
	splitquerylist := make([]string, 0)
	if strings.HasPrefix(querykey, JSONNTIMESKEY) == true {
		splitquerylist = strings.Split(querykey, fmt.Sprintf("%s%s", JSONNTIMESKEY, JSONSPLITKEY))
		splitquerylist[0] = JSONSPLITKEY
	} else {
		splitquerylist = strings.Split(querykey, fmt.Sprintf("%s%s%s", JSONSPLITKEY, JSONNTIMESKEY, JSONSPLITKEY))
	}
	jsonarray, err := jsondata.Query(splitquerylist[0])
	if err != nil {
		return nil, err
	}
	switch jsonarraydata := jsonarray.(type) {
	case map[string]interface{}:
		if valone, ok := jsonarraydata[splitquerylist[1]]; ok {
			return valone, nil
		}
	case []interface{}:
		datalist := make([]interface{}, 0)
		for _, jsonarrayval := range jsonarraydata {
			arrayjsonparser, err := iJd.NewJSONData(jsonarrayval)
			if err != nil {
				return nil, err
			}
			jsonval, err := arrayjsonparser.Query(splitquerylist[1])
			if err != nil {
				return nil, err
			}
			datalist = append(datalist, jsonval)
		}
		return datalist, nil
	}
	return nil, errors.New("No Convert")
}

// getRecordsetWithQueryKey get Recordset with query key
func (jc JSONConvert) getRecordsetWithQueryKey(querykey string, convertdata map[string]interface{}, jsonparser *iJd.JSONData) ([]interface{}, error) {
	if datasetquerykey, ok := convertdata[querykey].(string); ok {
		jsonval, err := jsonparser.Query(datasetquerykey)
		if err != nil {
			return nil, err
		}
		switch jsonvaldata := jsonval.(type) {
		case []interface{}:
			return jsonvaldata, nil
		case interface{}:
			return []interface{}{jsonvaldata}, nil
		}
	}
	return nil, fmt.Errorf("No RecordSet %s", querykey)
}

// getDataFromJoinRecordset get data from JoinRecordset
func (jc JSONConvert) getDataFromJoinRecordset(joincolumn string, joincolumnvalue interface{}, joindata []interface{}) interface{} {
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
