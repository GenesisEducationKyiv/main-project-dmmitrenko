package model

type Rate struct {
	Rates map[string]map[string]float64 `json:"rates"`
}
