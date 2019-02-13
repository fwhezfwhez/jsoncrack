package jsoncrack

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"math"
	"reflect"
	"time"
)

func SmartPrint(i interface{}, escapeZero ...bool) {
	if len(escapeZero) > 1 {
		panic(errorx.NewFromStringf("'escapeZero' should be length by 1 or 0 but got %d", len(escapeZero)))
	}

	var kv = make(map[string]interface{})
	vValue := reflect.ValueOf(i)
	vType := reflect.TypeOf(i)
	for i := 0; i < vValue.NumField(); i++ {
		kv[vType.Field(i).Name] = vValue.Field(i).Interface()
	}
	fmt.Println("receive:")
	for k, v := range kv {
		if len(escapeZero) == 1 && escapeZero[0] == true {
			if IfZero(v) {
				continue
			}
		}
		fmt.Print(k)
		fmt.Print(":")
		fmt.Print(v)
		fmt.Println()
	}
}

func IfZero(arg interface{}) bool {
	if arg == nil {
		return true
	}
	switch v := arg.(type) {
	case int, int32, int16, int64:
		if v == 0 {
			return true
		}
	case float32:
		r := float64(v)
		return math.Abs(r-0) < 0.0000001
	case float64:
		return math.Abs(v-0) < 0.0000001
	case string:
		if v == "" || v == "%%" || v == "%" {
			return true
		}
	case *string, *int, *int64, *int32, *int16, *int8, *float32, *float64, *time.Time:
		if v == nil {
			return true
		}
	case time.Time:
		return v.IsZero()
	case Time:
		return v.Time().IsZero()
	default:
		return false
	}
	return false
}
