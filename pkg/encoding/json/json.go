package json

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

func MarshalToString(v interface{}) string {
	b, err := Marshal(v)
	if err != nil {
		fmt.Println("json string" + err.Error())
		return ""
	}
	return string(b)
}
