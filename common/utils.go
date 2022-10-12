package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log-ext/common/errorx"
	"log-ext/common/response"
	"reflect"
	"strconv"
	"time"
)

func ParseBasicAuth(auth string) (key string, ok bool) {
	if auth == "" {
		return "", false
	}
	const prefix = "Bearer "
	if len(auth) < len(prefix) || !EqualFold(auth[:len(prefix)], prefix) {
		return "", false
	}
	return auth[len(prefix):], true
}

func EqualFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] != t[i] {
			return false
		}
	}
	return true
}

func ThrowErr(err *errorx.CodeError) response.HttpResponse {
	return response.HttpResponse{
		Code: uint32(err.GetErrCode()),
		Msg:  err.GetErrMsg(),
		Data: struct{}{},
	}
}

func StringValue(a *string) string {
	if a == nil {
		return ""
	}
	return *a
}

func ParseJSON(a *string) interface{} {
	mapTmp := make(map[string]interface{})
	d := json.NewDecoder(bytes.NewReader([]byte(StringValue(a))))
	d.UseNumber()
	err := d.Decode(&mapTmp)
	if err == nil {
		return mapTmp
	}

	sliceTmp := make([]interface{}, 0)
	d = json.NewDecoder(bytes.NewReader([]byte(StringValue(a))))
	d.UseNumber()
	err = d.Decode(&sliceTmp)
	if err == nil {
		return sliceTmp
	}

	if num, err := strconv.Atoi(StringValue(a)); err == nil {
		return num
	}

	if ok, err := strconv.ParseBool(StringValue(a)); err == nil {
		return ok
	}

	if floa64tVal, err := strconv.ParseFloat(StringValue(a), 64); err == nil {
		return floa64tVal
	}
	return nil
}

func AssertAsMap(a interface{}) map[string]interface{} {
	r := reflect.ValueOf(a)
	if r.Kind().String() != "map" {
		panic(fmt.Sprintf("%v is not a map[string]interface{}", a))
	}

	res := make(map[string]interface{})
	tmp := r.MapKeys()
	for _, key := range tmp {
		res[key.String()] = r.MapIndex(key).Interface()
	}

	return res
}

func String(a string) *string {
	return &a
}

func ReadAsString(body io.Reader) (*string, error) {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return String(""), err
	}
	r, ok := body.(io.ReadCloser)
	if ok {
		r.Close()
	}
	return String(string(byt)), nil
}

// 索引拼接
func StructIndexs(typ string, t1, t2 time.Time) (indexs []string) {
	type IndexStruc struct {
		Env     int `json:"env",required:"true"`
		Region  int `json:"region",required:"true"`
		Project int `json:"project",required:"true"`
	}

	ids := []IndexStruc{{Env: 1, Project: 1, Region: 1}, {Env: 2, Project: 2, Region: 2}}

	y1 := t1.Year()
	y2 := t2.Year()
	m1 := int(t1.Month())
	m2 := int(t2.Month())

	y, m := 0, 0
	done := false

	for y = y1; y <= y2; y++ {
		if done {
			break
		}
		item := typ + "-%v-%v-%v-%v-%v"

		if y == y1 && y1 < y2 {
			for m = m1; m <= 12; m++ {
				for _, v := range ids {
					indexs = append(indexs, fmt.Sprintf(item, v.Region, v.Project, v.Env, y, m))
				}
			}
		}

		if y == y2 && y1 < y2 {
			for m = 1; m <= m2; m++ {
				if m == m2 {
					done = true
					break
				}

				for _, v := range ids {
					indexs = append(indexs, fmt.Sprintf(item, v.Region, v.Project, v.Env, y, m))
				}
			}
			break
		}

		if y1 == y2 {
			for m = m1; m <= m2; m++ {

				for _, v := range ids {
					indexs = append(indexs, fmt.Sprintf(item, v.Region, v.Project, v.Env, y, m))
				}
			}
			break
		}

		for m = 0; m <= 12; m++ {
			for _, v := range ids {
				indexs = append(indexs, fmt.Sprintf(item, v.Region, v.Project, v.Env, y, m))
			}
		}

	}

	return indexs
}
