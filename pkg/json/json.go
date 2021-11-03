package json

import jsoniter "github.com/json-iterator/go"

// RawMessage ...
type RawMessage = jsoniter.RawMessage

// Replace methods of "encoding/json"
var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)
