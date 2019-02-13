package main

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/jsoncrack"
)

func main() {
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

	// init
	jc := jsoncrack.NewCracker(nil)
	// get
	fmt.Println("get class.master")
	r, e := jc.Get(jsoncrack.BYTES, in, "class", "master")
	if e != nil {
		panic(e)
	}
	fmt.Println(string(r.([]byte)))

	// update
	fmt.Println("update class.master.company.country.location")
	r, e = jc.Update(in, "location", "亚洲", "class", "master", "company", "country")
	if e != nil {
		panic(e)
	}
	fmt.Println(string(r.([]byte)))

	// add
	fmt.Println("add class.master.company.country.chinese_name : '中国'")
	r, e = jc.Add(in, "chinese_name", "中国", "class", "master", "company", "country")
	if e != nil {
		panic(e)
	}
	fmt.Println(string(r.([]byte)))

	// delete
	fmt.Println("delete class.master.company.manager")
	r, e = jc.Delete(jsoncrack.BYTES, false, in, "class", "master", "company", "manager")
	if e != nil {
		panic(e)
	}
	fmt.Println(string(r.([]byte)))

	type Time = jsoncrack.Time
	type VO struct {
		CreatedAt1  Time `json:"created_at1"`
		CreatedAt2  Time `json:"created_at2"`
		CreatedAt3  Time `json:"created_at3"`
		CreatedAt4  Time `json:"created_at4"`
		CreatedAt6  Time `json:"created_at6"`
		CreatedAt7  Time `json:"created_at7"`
		CreatedAt8  Time `json:"created_at8"`
		CreatedAt9  Time `json:"created_at9"`
		CreatedAt10 Time `json:"created_at10"`
		CreatedAt11 Time `json:"created_at11"`
		CreatedAt12 Time `json:"created_at12"`
		CreatedAt13 Time `json:"created_at13"`
	}
	var request = []byte(`{
        "created_at1": "2018-01-01",
        "created_at2": "2018-1-01",
		"created_at3": "2018/01/01",
		"created_at4": "2018/1/01",
        "created_at6": "2018.01.01",
        "created_at7": "2018.1.01",

        "created_at8": "2018-01-01 15:04:05",
        "created_at9": "2018-1-01 15:04:05",
		"created_at10": "2018.01.01 15:04:05",
		"created_at11": "2018.1.01 15:04:05",
		"created_at12": "2018/01/01 15:04:05",
		"created_at13": "2018/1/01 15:04:05"
    }`)
	vo := VO{}
	e = json.Unmarshal(request, &vo)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	jsoncrack.SmartPrint(vo, true)
}
