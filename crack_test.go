package jsoncrack

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewCracker(t *testing.T) {
	jc := NewCracker(nil)
	buf, e := jc.Json.Marshal(struct{ Name string }{Name: "test_cracker_(un)marshal"})
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
	fmt.Println(string(buf))
	var v map[string]interface{}

	e = jc.Json.Unmarshal(buf, &v)
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
	fmt.Println(v)
}

func TestGenRawJson(t *testing.T) {
	jc := NewCracker(nil)
	rs := string(jc.generateRawJson("name", "ft", "key1", "key2", "key3"))
	fmt.Println(rs)
	if rs != `{"key1":{"key2":{"key3":{"name":"ft"}}}}` {
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	var in = []byte(`{
		"class": {
			"name": "高中1班",
			"master": {
				"name": "张一山",
				"age": 21,
				"company": {
					"name": "go公司",
					"built_by": "张二山",
					"manager": ["张一山", "张二山", "张三山"],
					"country": {
						"name": "China",
						"location": "Asure"
					}
				}
			}
		}
	}`)

	jc := NewCracker(nil)

	// add a new field chinese_name
	b, e := jc.Update(in, "chinese_name", "中国", "class", "master", "company", "country")
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
	fmt.Println(string(b))

	// update an existed field chinese_name
	b, e = jc.Update(in, "location", "亚洲", "class", "master", "company", "country")
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
	fmt.Println(string(b))

	// safeUpdate throws error while exists key 'location'
	b, e = jc.SafeUpdate(in, "location", "亚洲", "class", "master", "company", "country")
	if e != nil {
		fmt.Println(e.Error())
	}
}

func TestGet(t *testing.T) {
	jc := NewCracker(nil)
	var in = []byte(`{
		"class": {
			"name": "高中1班",
			"master": {
				"name": "张一山",
				"age": 21,
				"company": {
					"name": "go公司",
					"built_by": "张二山",
					"manager": ["张一山", "张二山", "张三山"],
					"country": {
						"name": "China",
						"location": "Asure"
					}
				}
			}
		}
	}`)

	r, _ := jc.Get(BYTES, in, "HAPPY")
	//if e != nil {
	//	fmt.Println(e.Error())
	//	t.Fatal()
	//}

	fmt.Println(r==nil)
	fmt.Println(111,string(r.([]byte)))

	var e error
	r, e = jc.Get(BYTES, in, "class", "master", "company", "manager")
	if e != nil {
		fmt.Println(e.Error())
		t.Fatal()
	}
	fmt.Println(string(r.([]byte)), reflect.TypeOf(r))
}

func TestDelete(t *testing.T) {
	var in = []byte(`{
		"class": {
			"name": "高中1班",
			"master": {
				"name": "张一山",
				"age": 21,
				"company": {
					"name": "go公司",
					"built_by": "张二山",
					"manager": ["张一山", "张二山", "张三山"],
					"country": {
						"name": "China",
						"location": "Asure"
					}
				}
			}
		}
	}`)

	jc := NewCracker(nil)

	r, e := jc.Delete(BYTES, false, in, "class", "name")
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
	fmt.Println(string(r.([]byte)))
}

func TestGetString(t *testing.T) {
	var in = []byte(`{
		"class": {
			"name": "高中1班",
			"master": {
				"name": "张一山",
				"age": 21,
				"company": {
					"name": "go公司",
					"built_by": "张二山",
					"manager": ["张一山", "张二山", "张三山"],
					"country": {
						"name": "China",
						"location": "Asure"
					}
				}
			}
		}
	}`)

	fmt.Println(GetString(in, "class", "master", "company", "name"))
}
