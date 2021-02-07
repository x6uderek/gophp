package gophp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	resp, err := http.Get("http://localhost/serialize.php")
	if err != nil {
		t.Fatal(err)
	}
	testStr,err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	tests := strings.Split(string(testStr),"\n")

	resp1, err := http.Get("http://localhost/serialize.php?q=json")
	if err != nil {
		t.Fatal(err)
	}
	expectStr,err := ioutil.ReadAll(resp1.Body)
	defer resp1.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	expects := strings.Split(string(expectStr), "\n")

	for i:=0; i<len(tests); i++ {
		parsed,_,err := Parse([]byte(tests[i]))
		if err != nil {
			t.Errorf("parse %s got err: %s", tests[i], err)
			continue
		}
		jsonStr,_ := json.Marshal(parsed)
		//go json用 \uaaaa, php用\uAAAA
		//go map的遍历顺序不确定，json输出也会变化，此处测试可能会失败
		if strings.ToUpper(string(jsonStr)) != strings.ToUpper(expects[i]) {
			t.Errorf("parse %s got %s but expect %s", tests[i], string(jsonStr), expects[i])
		}
	}
}