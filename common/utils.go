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
