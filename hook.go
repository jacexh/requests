package requests

import (
	"encoding/json"
	"encoding/xml"
)

func UnmarshalJSON(v any) Unmarshaller {
	return func(data []byte) error {
		return json.Unmarshal(data, v)
	}
}

func UnmarlXML(v any) Unmarshaller {
	return func(data []byte) error {
		return xml.Unmarshal(data, v)
	}
}
