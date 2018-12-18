package main

import (
	"fmt"
	"jsoncrack"
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
	r ,e :=jc.Get(jsoncrack.BYTES, in, "class", "master")
	if e!=nil {
		panic(e)
	}
	fmt.Println(string(r.([]byte)))

	// update
	fmt.Println("update class.master.company.country.location")
	r, e = jc.Update(in, "location", "亚洲", "class", "master", "company", "country")
	if e!=nil {
		panic(e)
	}
	fmt.Println(string(r.([]byte)))

	// add
	fmt.Println("add class.master.company.country.chinese_name : '中国'")
	r, e = jc.Add(in, "chinese_name", "中国", "class", "master", "company", "country")
	if e!=nil {
		panic(e)
	}
	fmt.Println(string(r.([]byte)))

	// delete
	fmt.Println("delete class.master.company.manager")
	r, e = jc.Delete(jsoncrack.BYTES,false,in, "class", "master", "company", "manager")
	if e!=nil {
		panic(e)
	}
	fmt.Println(string(r.([]byte)))
}
