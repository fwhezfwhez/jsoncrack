package jsoncrack

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"reflect"
	"strings"
)

const(
	MAP = 1+iota
	ARRAY
	BYTES
)

type JsonMarshaller interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}
type JsonCracker struct {
	Json JsonMarshaller
}

// a jsonMarshaller realization via the officials
type Jsoner struct {
}

func (j Jsoner) Marshal(dest interface{}) ([]byte, error) {
	return json.Marshal(dest)
}
func (j Jsoner) Unmarshal(data []byte, dest interface{}) error {
	return json.Unmarshal(data, dest)
}

// json marshal via its jsoner realization
func (jc JsonCracker) Marshal(dest interface{}) ([]byte, error) {
	return jc.Json.Marshal(dest)
}

// json unmarshal via its jsoner realization
func (jc JsonCracker) Unmarshal(data []byte, dest interface{}) error {
	return jc.Json.Unmarshal(data, dest)
}

// New a json-cracker instance
// jsoner == nil , official json marshaler will be put to use,otherwise use jsoner as specific
func NewCracker(jsoner JsonMarshaller) JsonCracker {
	if jsoner == nil {
		jsoner = Jsoner{}
	}
	return JsonCracker{
		Json: jsoner,
	}
}

// update a json raw-message with a new key-value pair
// assume data:
//{
//    "class1": {
//        "master": {
//            "name": "Li Hua",
//            "Age": 28
//        },
//       "students": [
//           {
//            "name": "Li Lei",
//            "Age": 12
//           },
//           {
//            "name": "Tom",
//            "Age": 11
//           }
//       ]
//    }
//}
// after exec jc.Update([]byte(data), "sub_name", "Li li", "class1", "master")
// returns:
//{
//    "class1": {
//        "master": {
//            "name": "Li Hua",
//            "Age": 28,
//            "sub_name":"Li li"
//        },
//       "students": [
//           {
//            "name": "Li Lei",
//            "Age": 12
//           },
//           {
//            "name": "Tom",
//            "Age": 11
//           }
//       ]
//    }
//}
func (jc JsonCracker) Update(data []byte, k string, v interface{}, keys ...string) ([]byte, error) {
	return jc.MustUpdate(data, k, v, keys...)
}

// add k-v pair into a json []byte
func (jc JsonCracker) Add(data []byte, k string, v interface{}, keys ...string) ([]byte, error){
	return jc.Update(data , k , v , keys ...)
}
// When exists the key 'k' , it replace the origin value with v
func (jc JsonCracker) MustUpdate(data []byte, k string, v interface{}, keys ...string) ([]byte, error) {
	var dest map[string]interface{}
	var err error
	if err = jc.Unmarshal(data, &dest); err != nil {
		return nil, errorx.NewFromString(fmt.Sprintf("data is not well format as json serial,error: '%s'", err.Error()))
	}
	keys = append(keys, k)

	if err = updateParse(false, dest, v, keys...); err != nil {
		return nil, errorx.NewFromString(fmt.Sprintf("error: '%s'", err.Error()))
	}
	return jc.Marshal(dest)
}

// when exists the key 'k', it throws an error without replacing
func (jc JsonCracker) SafeUpdate(data []byte, k string, v interface{}, keys ...string) ([]byte, error) {
	var dest map[string]interface{}
	var err error
	if err = jc.Unmarshal(data, &dest); err != nil {
		return nil, errorx.NewFromString(fmt.Sprintf("data is not well format as json serial,error: '%s'", err.Error()))
	}
	keys = append(keys, k)

	if err = updateParse(true, dest, v, keys...); err != nil {
		return nil, errorx.NewFromString(fmt.Sprintf("error: '%s'", err.Error()))
	}
	return jc.Marshal(dest)
}

// get a value via keys through a json []byte.
// value can specifc type as jsoncrack.BYTES or jsoncrack.MAP.The formmer returns json []bytes boxing in interface{},the
// other is a map[string]interface{} boxing in interface{}
func (jc JsonCracker) Get(vtype int, data []byte, keys ...string) (interface{}, error) {
	var dest map[string]interface{}
	er := jc.Unmarshal(data, &dest)
	if er!=nil {
		return nil, errorx.NewFromString(fmt.Sprintf("'data' is not well json format,err: %s", er.Error()))
	}
	var result interface{}
	er = getParse(&result, dest, keys ...)
	if er!=nil {
		return nil, errorx.NewFromString(fmt.Sprintf("get parse err: %s", er.Error()))
	}
	if vtype == MAP {
		return result,nil
	} else{
		return jc.Marshal(result)
	}
}

// vtype: jsoncrack.MAP,jsoncrack.Array,jsoncrack.BYTES
// when safe is true, delete the specific key while not existed, throws an error.
// when safe is false, do nothing if not existed keys
func (jc JsonCracker) Delete(vtype int,safe bool,data []byte, keys ...string)(interface{}, error){
	var dest map[string]interface{}
	er := jc.Unmarshal(data, &dest)
	if er!=nil {
		return nil, errorx.NewFromString(fmt.Sprintf("'data' is not well json format,err: %s", er.Error()))
	}
	er = deleteParse(safe, dest, keys ...)
	if er!=nil {
		return nil, errorx.NewFromString(fmt.Sprintf("get parse err: %s", er.Error()))
	}
	if vtype == MAP {
		return dest,nil
	} else{
		return jc.Marshal(dest)
	}
}

func (jc JsonCracker) generateRawJson(k string, v interface{}, keySrials ...string) []byte {
	var iterator = make(map[string]interface{}, 0)
	var empty = make(map[string]interface{}, 0)
	for i := len(keySrials) - 1; i >= 0; i-- {
		var tmp = CopyFrom(empty)

		if i == len(keySrials)-1 {
			iterator[k] = v
		}
		tmp[keySrials[i]] = CopyFrom(iterator)
		iterator = CopyFrom(tmp)
		if i == 0 {
			buf, _ := jc.Marshal(tmp)
			return buf
		}
	}
	return nil
}

// copy a map to another
func CopyFrom(m map[string]interface{}) map[string]interface{} {
	var m2 = make(map[string]interface{}, 0)
	for k, v := range m {
		m2[k] = v
	}
	return m2
}

// parse the 'target' typed map/struct via 'keys' and add/update 'val' into the target.
// while 'safe' is true, it throws error when the last key of keys exist, otherwise replace it if existed
func updateParse(safe bool, target interface{}, val interface{}, keys ...string) error {
	return updateParseDetail(safe, reflect.ValueOf(&target), keys, reflect.ValueOf(val))
}

// detail of parse
func updateParseDetail(safe bool, target reflect.Value, sel []string, val reflect.Value) error {
	if len(sel) == 0 {
		return errorx.NewFromString("invalid sel")
	}
	v := reflect.Indirect(target)
	//v := reflect.ValueOf(target)

L:
	switch v.Kind() {
	case reflect.Struct:
		nam := strings.Title(sel[0])
		if len(sel) == 1 {
			f := v.FieldByName(nam)
			if f.IsValid() {
				if f.CanSet() {
					f.Set(val)
					return nil
				}
			}
		} else {
			f := v.FieldByName(nam)
			return updateParseDetail(safe, f, sel[1:], val)
		}
	case reflect.Map:
		// nam := reflect.ValueOf(strings.Title(sel[0]))
		nam := reflect.ValueOf(sel[0])
		if len(sel) == 1 {
			//v.SetMapIndex(nam, val)
			vvvv := v.Interface().(map[string]interface{})
			if !safe {
				vvvv[nam.String()] = val.Interface()
			} else {
				if _, ok := vvvv[nam.String()]; ok {
					return errorx.NewFromString(fmt.Sprintf("key '%s' already exist", nam.String()))
				}
			}
			//f := v.MapIndex(nam)
			//if f.IsValid() {
			//	if f.CanSet() {
			//		f.Set(val)
			//		return nil
			//	} else {
			//		f = f.Elem()
			//		f.Set(val)
			//		return nil
			//	}
			//} else {
			//	v.SetMapIndex(nam, val)
			//}
		} else {
			f := v.MapIndex(nam)
			return updateParseDetail(safe, f, sel[1:], val)
		}
	case reflect.Interface:
		vv := v.Interface()
		if vvv, ok := vv.(map[string]interface{}); ok {
			v = reflect.ValueOf(vvv)
			goto L
		}
		fallthrough
	default:
		return errorx.NewFromString("invalid type, must struct/map")
	}
	return nil
}

// get the value from 'target', typed map/struct, indexed by keys...
func getParse(result *interface{},target interface{}, keys...string)error{
	return getParseDetail(result, reflect.ValueOf(&target), keys...)
}
// detail of getParse
func getParseDetail(result *interface{},target reflect.Value, sel ...string)error{
	v:= reflect.Indirect(target)
	if len(sel)==0 {
		return nil
	}

L:
	switch v.Kind() {
	case reflect.Array,reflect.Slice:
		nam := strings.Title(sel[0])
		if len(sel) ==1 {
			f := v.FieldByName(nam)
			*result = f.Interface()
		} else{
			return errorx.NewFromString(fmt.Sprintf("key '%s' is a array/slice,however keys is not end", nam))
		}
	case reflect.Struct:
		nam := strings.Title(sel[0])
		if len(sel) == 1 {
			f := v.FieldByName(nam)
			*result = f.Interface()
		} else {
			f := v.FieldByName(nam)
			return getParseDetail(result, f, sel[1:]...)
		}
	case reflect.Map:
		// nam := reflect.ValueOf(strings.Title(sel[0]))
		nam := reflect.ValueOf(sel[0])
		if len(sel) == 1 {
			*result = v.Interface().(map[string]interface{})[sel[0]]
		} else {
			f := v.MapIndex(nam)
			return getParseDetail(result, f, sel[1:]...)
		}
	case reflect.Interface:
		vv := v.Interface()
		if vvv, ok := vv.(map[string]interface{}); ok {
			_,ok2 :=vvv[sel[0]]
			if !ok2{
				return errorx.NewFromString(fmt.Sprintf("key '%s' not found", sel[0]))
			}
			v = reflect.ValueOf(vvv)
			goto L
		}
		fallthrough
	default:
		return errorx.NewFromString("invalid type, must struct/map")
	}
	return nil
}

// delete a key-value pair of a target via keys indexing
func deleteParse(safe bool, target interface{}, keys ...string)error{
	return deleteParseDetail(safe, reflect.ValueOf(&target), keys ...)
}

func deleteParseDetail(safe bool, target reflect.Value, keys...string)error{
	if len(keys) == 0 {
		return errorx.NewFromString("invalid sel")
	}
	v := reflect.Indirect(target)
	//v := reflect.ValueOf(target)

L:
	switch v.Kind() {
	case reflect.Map:
		// nam := reflect.ValueOf(strings.Title(sel[0]))
		nam := reflect.ValueOf(keys[0])
		if len(keys) == 1 {
			//v.SetMapIndex(nam, val)
			vvvv := v.Interface().(map[string]interface{})
			if !safe {
				delete(vvvv, nam.String())
			} else {
				if _, ok := vvvv[nam.String()]; !ok {
					return errorx.NewFromString(fmt.Sprintf("key '%s' not exist", nam.String()))
				}
			}
		} else {
			f := v.MapIndex(nam)
			return deleteParseDetail(safe, f, keys[1:]...)
		}
	case reflect.Interface:
		vv := v.Interface()
		if vvv, ok := vv.(map[string]interface{}); ok {
			v = reflect.ValueOf(vvv)
			goto L
		}
		fallthrough
	default:
		return errorx.NewFromString("invalid type, must struct/map")
	}
	return nil
}
