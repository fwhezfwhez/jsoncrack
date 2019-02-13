# jsoncrack
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/fwhezfwhez/jsoncrack)
[![Build Status]( https://www.travis-ci.org/fwhezfwhez/jsoncrack.svg?branch=master)]( https://www.travis-ci.org/fwhezfwhez/jsoncrack)

jsoncrack is tool on developing to operate json []byte straightly.
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [1. Start](#1-start)
- [2. Notice](#2-notice)
  - [2.1 json marshaller can be specific when inited.](#21-json-marshaller-can-be-specific-when-inited)
  - [2.2 crud of jsoncrack doesn't change the original data.](#22-crud-of-jsoncrack-doesnt-change-the-original-data)
  - [2.3 Get() and Delete() can specific the returning type of the modified copy of data,ranging in `[jsoncrack.MAP, jsoncrack.ARRAY, jsoncrack.BYTE]`](#23-get-and-delete-can-specific-the-returning-type-of-the-modified-copy-of-dataranging-in-jsoncrackmap-jsoncrackarray-jsoncrackbyte)
- [3. Function menus](#3-function-menus)
- [4. Example](#4-example)
- [5. Jsoncrack.Time](#5-jsoncracktime)
  - [5.1 Available layouts](#51-available-layouts)
  - [5.2 Example](#52-example)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 1. Start
`go get github.com/fwhezfwhez/jsoncrack`

## 2. Notice
### 2.1 json marshaller can be specific when inited.
Any instance that realize JsonMarshaller interface can be set jsoncrack's json marshaller. The default is the official json package `encoding/json`.

```go
type JsonMarshaller interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type Jsoner struct {
}

func (j Jsoner) Marshal(dest interface{}) ([]byte, error) {
	return json.Marshal(dest)
}

func (j Jsoner) Unmarshal(data []byte, dest interface{}) error {
	return json.Unmarshal(data, dest)
}
func main() {
    jsonMarshaler := Jsoner{}
	jc := jsoncrack.NewCracker(jsonMarshaler)
	....
}

```
### 2.2 crud of jsoncrack doesn't change the original data.
All apis of jsoncrack which operates []byte data returns a modified copy of the original data.

### 2.3 Get() and Delete() can specific the returning type of the modified copy of data,ranging in `[jsoncrack.MAP, jsoncrack.ARRAY, jsoncrack.BYTE]`
jsoncrack.Bytes is ok for all cases.If you ensure your returning data is a formated map[string]interface{}, json.Map is optional.
jsoncrack.ARRAY is now the same as BYTES.It might be changed with upgrading the version of jsoncrack.

## 3. Function menus
| function | info |
|:----------- | :---- |
| Marshal | json encoding  |
| Unmarshal | json decoding  |
| Update | update json []byte field  |
| Delete | delete json []byte field|
| Add | add json []byte field|
| Get | get json []byte field value |

## 4. Example
```go
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
```

## 5. Jsoncrack.Time
### 5.1 Available layouts
jsoncrack.Time can receive all kinds of time string format layouts below:
```go
		"2006-01-02",
		"2006-1-2",

		"2006/01/02",
		"2006/1/2",

		"2006.01.02",
		"2006.1.2",

		"2006-01-02 15:04:05",
		"2006-1-2 15:04:05",

		"2006/01/02 15:04:05",
		"2006/1/2 15:04:05",

		"2006.01.02 15:04:05",
		"2006.1.2 15:04:05",
```
### 5.2 Example
```go
type Time jsoncrack.Time
func main(){
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
	e := json.Unmarshal(request, &vo)
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
		return
	}
	jsoncrack.SmartPrint(vo, true)
}
```
