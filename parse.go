package gophp

import (
	"fmt"
	"regexp"
	"strconv"
)

func Parse(source []byte) (interface{}, int, error) {
	for i := 0; i < len(source); {
		switch source[i] {
		case 'b':
			boolValue, length, err := parseBoolean(source[i:])
			if err != nil {
				return nil, 0, fmt.Errorf("%s at %d", err.Error(), i)
			}
			i += length
			return boolValue, length, nil
		case 'i':
			intValue, length, err := parseInt(source[i:])
			if err != nil {
				return nil, 0, fmt.Errorf("%s at %d", err.Error(), i)
			}
			i += length
			return intValue, length, nil
		case 'd':
			floatValue, length, err := parseFloat(source[i:])
			if err != nil {
				return nil, 0, fmt.Errorf("%s at %d", err.Error(), i)
			}
			i += length
			return floatValue, length, nil
		case 's':
			stringValue, length, err := parseString(source[i:])
			if err != nil {
				return nil, 0, fmt.Errorf("%s at %d", err.Error(), i)
			}
			i += length
			return stringValue, length, nil
		case 'a':
			arrayValue, length, err := parseArray(source[i:])
			if err != nil {
				return nil, 0, fmt.Errorf("%s at %d", err.Error(), i)
			}
			i += length
			return arrayValue, length, nil
		case 'N':
			nilValue, length, err := parseNil(source[i:])
			if err != nil {
				return nil, 0, fmt.Errorf("%s at %d", err.Error(), i)
			}
			i += length
			return nilValue, length, nil
		case 'O':
			objectValue, length, err := parseObject(source[i:])
			if err != nil {
				return nil, 0, fmt.Errorf("%s at %d", err.Error(), i)
			}
			i += length
			return objectValue, length, nil
		default:
			return nil, 0, fmt.Errorf("unsupported type at %d", i)
		}
	}
	return nil, 0, fmt.Errorf("invalid params")
}

func parseBoolean(source []byte) (bool, int, error) {
	if err := shouldBeginWith(source, 'b', 4); err != nil {
		return false, 0, err
	}
	var value bool
	switch source[2] {
	case '1':
		value = true
	case '0':
		value = false
	default:
		return false, 0, fmt.Errorf("should get '0' or '1' when parse bool")
	}
	if err := shouldEnd(source, 'b', 3); err != nil {
		return false, 0, err
	}
	return value, 4, nil
}

func parseInt(source []byte) (int64, int, error) {
	if err := shouldBeginWith(source, 'i', 4); err != nil {
		return 0, 0, err
	}
	i := 2
	for ; source[i] >= '0' && source[i] <= '9' || source[i] == '-'; i++ {
		//nothing here
	}
	if err := shouldEnd(source, 'i', i); err != nil {
		return 0, 0, err
	}
	stringValue := string(source[2:i])
	exp := regexp.MustCompile(`^[-+]?\d+$`)
	if !exp.MatchString(stringValue) {
		return 0, 0, fmt.Errorf("invalid int value: %s", stringValue)
	}
	value, _ := strconv.ParseInt(stringValue, 10, 64) //ignore error
	return value, i + 1, nil
}

func parseFloat(source []byte) (float64, int, error) {
	if err := shouldBeginWith(source, 'd', 5); err != nil {
		return 0, 0, err
	}
	i := 2
	for ; source[i] >= '0' && source[i] <= '9' || source[i] == '-' || source[i] == '.'; i++ {
		//nothing here
	}
	if err := shouldEnd(source, 'd', i); err != nil {
		return 0, 0, err
	}
	stringValue := string(source[2:i])
	exp := regexp.MustCompile(`^[-+]?\d+(\.\d*)?$`)
	if !exp.MatchString(stringValue) {
		return 0, 0, fmt.Errorf("invalid float value: %s", stringValue)
	}
	value, _ := strconv.ParseFloat(stringValue, 64)
	return value, i + 1, nil
}

func parseString(source []byte) (string, int, error) {
	if err := shouldBeginWith(source, 's', 8); err != nil {
		return "", 0, err
	}
	i := 2
	//获取string长度
	for ; source[i] >= '0' && source[i] <= '9'; i++ {
		//nothing
	}
	stringLength, _ := strconv.ParseInt(string(source[2:i]), 10, 64)
	if source[i] != ':' {
		return "", 0, fmt.Errorf("should get ':' when parse string")
	}
	i++
	if source[i] != '"' {
		return "", 0, fmt.Errorf("should get '\"' when parse string")
	}
	if len(source) <= i+int(stringLength)+2 { //后面还有";
		return "", 0, fmt.Errorf("invalid string length")
	}
	i++
	value := string(source[i : i+int(stringLength)])
	i += int(stringLength)
	if source[i] != '"' {
		return "", 0, fmt.Errorf("should get '\"' when parse string")
	}
	i++
	if err := shouldEnd(source, 's', i); err != nil {
		return "", 0, err
	}
	return value, i + 1, nil
}

//parseArray 可能返回切片或者map
//因为php的数组是可以混合的类型
func parseArray(source []byte) (interface{}, int, error) {
	if err := shouldBeginWith(source, 'a', 6); err != nil {
		return nil, 0, err
	}
	i := 2
	//获取array长度
	for ; source[i] >= '0' && source[i] <= '9'; i++ {
		//nothing
	}
	arrayLength, _ := strconv.ParseInt(string(source[2:i]), 10, 64)
	if source[i] != ':' {
		return nil, 0, fmt.Errorf("should get ':' when parse array")
	}
	i++
	if source[i] != '{' {
		return nil, 0, fmt.Errorf("should get '{' when parse array")
	}
	i++
	isPureArray := true //如果是纯净的数组，返回切片
	lastIntIndex := int64(-1)
	arrayValue := make([]interface{}, 0, arrayLength)  //切片返回值
	value := make(map[string]interface{}, arrayLength) //map返回值
	for index := 0; index < int(arrayLength); index++ {
		switch source[i] {
		case 'i':
			intKey, length, err := parseInt(source[i:])
			if err != nil {
				return nil, 0, err
			}
			i += length
			curValue, length, err := Parse(source[i:])
			if err != nil {
				return nil, 0, err
			}
			i += length
			value[fmt.Sprint(intKey)] = curValue
			if intKey == lastIntIndex+1 {
				lastIntIndex = intKey
				arrayValue = append(arrayValue, curValue)
			} else {
				isPureArray = false
			}
		case 's':
			stringKey, length, err := parseString(source[i:])
			if err != nil {
				return nil, 0, err
			}
			i += length
			curValue, length, err := Parse(source[i:])
			if err != nil {
				return nil, 0, err
			}
			i += length
			value[stringKey] = curValue
			isPureArray = false
		default:
			return nil, 0, fmt.Errorf("unsupported array key")
		}
	}
	if source[i] != '}' {
		return nil, 0, fmt.Errorf("should get '}' when parse array")
	}
	if isPureArray {
		return arrayValue, i + 1, nil
	}
	return value, i + 1, nil
}

func parseNil(source []byte) (interface{}, int, error) {
	if len(source) < 2 || source[0] != 'N' || source[1] != ';' {
		return nil, 0, fmt.Errorf("invalid nil value")
	}
	return nil, 2, nil
}

func parseObject(source []byte) (map[string]interface{}, int, error) {
	if err := shouldBeginWith(source, 'O', 12); err != nil {
		return nil, 0, err
	}
	i := 2
	//获取class name长度
	for ; source[i] >= '0' && source[i] <= '9'; i++ {
		//nothing
	}
	classNameLength, _ := strconv.ParseInt(string(source[2:i]), 10, 64)
	if source[i] != ':' {
		return nil, 0, fmt.Errorf("should get ':' when parse object")
	}
	i++
	if source[i] != '"' {
		return nil, 0, fmt.Errorf("should get '\"' when parse object")
	}
	if len(source) <= i+int(classNameLength)+6 { //后面还有":0:{}
		return nil, 0, fmt.Errorf("invalid class name length")
	}
	i++
	i += int(classNameLength) //ignore class name
	if source[i] != '"' {
		return nil, 0, fmt.Errorf("should get '\"' when parse object")
	}
	i++
	if source[i] != ':' {
		return nil, 0, fmt.Errorf("should get ':' when parse object")
	}
	i++
	curIndex := i
	//获取object长度
	for ; source[i] >= '0' && source[i] <= '9'; i++ {
		//nothing
	}
	objectLength, _ := strconv.ParseInt(string(source[curIndex:i]), 10, 64)
	if source[i] != ':' {
		return nil, 0, fmt.Errorf("should get ':' when parse array")
	}
	i++
	if source[i] != '{' {
		return nil, 0, fmt.Errorf("should get '{' when parse array")
	}
	i++
	value := make(map[string]interface{}, objectLength)
	for index := 0; index < int(objectLength); index++ {
		switch source[i] {
		case 's':
			stringKey, length, err := parseString(source[i:])
			if err != nil {
				return nil, 0, err
			}
			i += length
			curValue, length, err := Parse(source[i:])
			if err != nil {
				return nil, 0, err
			}
			i += length
			value[stringKey] = curValue
		default:
			return nil, 0, fmt.Errorf("unsupported object key")
		}
	}
	if source[i] != '}' {
		return nil, 0, fmt.Errorf("should get '}' when parse object")
	}
	return value, i + 1, nil
}

func shouldBeginWith(source []byte, t byte, minLen int) error {
	if len(source) < minLen {
		return fmt.Errorf("invalid int length")
	}
	if source[0] != t {
		return fmt.Errorf("should get '%c' when parse bool", t)
	}
	if source[1] != ':' {
		return fmt.Errorf("should have ':' when parse '%c'", t)
	}
	return nil
}

func shouldEnd(source []byte, t byte, index int) error {
	if len(source) <= index || source[index] != ';' {
		return fmt.Errorf("not end with ';' when parse '%c'", t)
	}
	return nil
}
