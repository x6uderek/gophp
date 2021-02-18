package gophp

import (
	"fmt"
	"math"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []string{
		`b:1;`,
		`b:0;`,
		`i:123;`,
		`d:3.1415;`,
		`s:43:"a ba a ba, "(){}?:"><?!@#$%^&*-=_+|'\中国";`,
		`N;`,
		`O:3:"cls":1:{s:5:"prop1";s:3:"ppp";}`,
		`a:4:{i:0;b:1;i:1;b:0;i:2;i:123;i:3;d:3.1415;}`,
		`a:7:{i:0;b:1;i:3;b:0;i:4;i:123;i:5;d:3.1415;s:3:"str";s:43:"a ba a ba, "(){}?:"><?!@#$%^&*-=_+|'\中国";s:4:"null";N;s:3:"obj";O:3:"cls":1:{s:5:"prop1";s:3:"ppp";}}`,
	}

	boolVal, _, err := Parse([]byte(testCases[0]))
	if err != nil {
		t.Errorf("parse %s got error: %v", testCases[0], err)
	}
	boolVal1 := boolVal.(bool)
	if boolVal1 != true {
		t.Errorf("parse %s got %v", testCases[0], boolVal1)
	}

	boolVal, _, err = Parse([]byte(testCases[1]))
	if err != nil {
		t.Errorf("parse %s got error: %v", testCases[1], err)
	}
	boolVal0 := boolVal.(bool)
	if boolVal0 != false {
		t.Errorf("parse %s got %v", testCases[1], boolVal0)
	}

	intVal, _, err := Parse([]byte(testCases[2]))
	if err != nil {
		t.Errorf("parse %s got error: %v", testCases[2], err)
	}
	int64Val := intVal.(int64)
	if int64Val != 123 {
		t.Errorf("parse %s got (%[2]T) %[2]v", testCases[2], int64Val)
	}

	floatVal, _, err := Parse([]byte(testCases[3]))
	if err != nil {
		t.Errorf("parse %s got error: %v", testCases[3], err)
	}
	float64Val := floatVal.(float64)
	//浮点数比较
	if math.Abs(float64Val - 3.1415) > 0.0001 {
		t.Errorf("parse %s got %v", testCases[3], float64Val)
	}

	stringVal, _, err := Parse([]byte(testCases[4]))
	if err != nil {
		t.Errorf("parse %s got error: %v", testCases[4], err)
	}
	stringStrVal := stringVal.(string)
	if stringStrVal != `a ba a ba, "(){}?:"><?!@#$%^&*-=_+|'\中国` {
		t.Errorf("parse %s got %v", testCases[4], stringStrVal)
	}

	nilVal, _, err := Parse([]byte(testCases[5]))
	if err != nil {
		t.Errorf("parse %s got error: %v", testCases[5], err)
	}
	if nilVal != nil {
		t.Errorf("parse %s got (%[2]T) %[2]v", testCases[5], nilVal)
	}

	objectVal, _, err := Parse([]byte(testCases[6]))
	if err != nil {
		t.Errorf("parse %s got error: %v", testCases[6], err)
	}
	objectObjVal := objectVal.(map[string]interface{})
	prop, ok := objectObjVal["prop1"]
	if !ok {
		t.Errorf("parse %s got no property: %v", testCases[6], "prop1")
	} else if propStr := prop.(string); propStr != "ppp" {
		t.Errorf("parse %s got property: %v value: %v", testCases[6], "prop1", propStr)
	}

	arrayVal, _, err := Parse([]byte(testCases[7]))
	if err != nil {
		t.Errorf("parse %s got error: %v", testCases[7], err)
	}
	arraySliceVal := arrayVal.([]interface{})
	if len(arraySliceVal) != 4 {
		t.Errorf("parse %s got array length: %d", testCases[7], len(arraySliceVal))
	}
	if item3 := arraySliceVal[2].(int64); item3 != 123 {
		t.Errorf("parse %s got item3: %d", testCases[7], item3)
	}

	fmt.Println(arraySliceVal)

	mapVal, _, err := Parse([]byte(testCases[8]))
	if err != nil {
		t.Errorf("parse %s got error: %v", testCases[8], err)
	}
	mapMVal := mapVal.(map[string]interface{})
	if len(mapMVal) != 7 {
		t.Errorf("parse %s got map length: %d", testCases[8], len(mapMVal))
	} else if itemStr, ok := mapMVal["str"]; !ok {
		t.Errorf("parse %s got no value for key: %s", testCases[8], "str")
	} else if itemString := itemStr.(string); itemString != `a ba a ba, "(){}?:"><?!@#$%^&*-=_+|'\中国` {
		t.Errorf("parse %s got %s value %v", testCases[8], "str", itemString)
	}

	fmt.Println(mapMVal)
}