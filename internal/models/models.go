package models

type Report struct {
	Id    string `json:"id" mongo:"_id"`
	Name  string `json:"name,omitempty" mongo:"name"`
	Query string `json:"query,omitempty" mongo:"query"`
	Data  string `json:"data,omitempty"`
}

type Diver struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Port    string `json:"port"`
	Michman string `json:"michman"`
}
