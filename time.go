package jsoncrack

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"time"
)

type Time time.Time
func (t Time) String()string{
	return fmt.Sprintf("%v", time.Time(t))
}
func (t Time) MarshalJSON()([]byte,error){
	return json.Marshal(time.Time(t))
}
func (t *Time) UnmarshalJSON(buf []byte) error {
	var e error
	var tStr string
	var t0 time.Time
	var errors = make([]error, 0, 6)

	e = json.Unmarshal(buf, &t0)
	if e == nil {
		*t = Time(t0)
		return nil
	}
	errors = append(errors, errorx.Wrap(e))
	e = json.Unmarshal(buf, &tStr)
	if e != nil {
		return e
	}

	for _, layout := range []string{
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
	} {
		t0, e = time.Parse(layout, tStr)
		if e == nil {
			*t = Time(t0)
			return nil
		}
		errors = append(errors, errorx.Wrap(e))
	}

	return errorx.GroupErrors(errors...)
}

func (t Time) Time() time.Time{
	return time.Time(t)
}
