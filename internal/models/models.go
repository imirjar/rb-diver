package models

import "encoding/json"

type Report struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Query string `json:"query,omitempty"`
	Data  string `json:"data,omitempty"`
}

type Data struct {
	Raw []json.RawMessage
}
