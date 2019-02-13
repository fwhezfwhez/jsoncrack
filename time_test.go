package jsoncrack

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTime(t *testing.T) {
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
	SmartPrint(vo, true)
}
