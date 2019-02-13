package jsoncrack

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/jsoncrack"
	"testing"
)

func TestTime(t *testing.T) {
	type VO struct {
		CreatedAt1  jsoncrack.Time `json:"created_at1"`
		CreatedAt2  jsoncrack.Time `json:"created_at2"`
		CreatedAt3  jsoncrack.Time `json:"created_at3"`
		CreatedAt4  jsoncrack.Time `json:"created_at4"`
		CreatedAt6  jsoncrack.Time `json:"created_at6"`
		CreatedAt7  jsoncrack.Time `json:"created_at7"`
		CreatedAt8  jsoncrack.Time `json:"created_at8"`
		CreatedAt9  jsoncrack.Time `json:"created_at9"`
		CreatedAt10 jsoncrack.Time `json:"created_at10"`
		CreatedAt11 jsoncrack.Time `json:"created_at11"`
		CreatedAt12 jsoncrack.Time `json:"created_at12"`
		CreatedAt13 jsoncrack.Time `json:"created_at13"`
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
