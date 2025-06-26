package models

type Report struct {
	Id          string `json:"id" bson:"_id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	Query string `json:"query,omitempty"`

	Roles []Role `json:"roles,omitempty"`

	Created string `json:"created,omitempty"`
	Updated string `json:"updated,omitempty"`
}

type Data struct {
	ID      string   `json:"id"`
	Columns []string `json:"columns"`
	Values  [][]any  `json:"values"`
}
