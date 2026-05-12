package eob

type Meta struct {
	LastUpdated string   `json:"lastUpdated"`
	Profile     []string `json:"profile"`
}

type Link struct {
	Relation string `json:"relation"`
	URL      string `json:"url"`
}

type Coding struct {
	System  string `json:"system"`
	Code    string `json:"code"`
	Display string `json:"display"`
}

type Type struct {
	Codings []Coding `json:"coding"`
}

type SubType struct {
	Codings []Coding `json:"coding"`
	Text    string   `json:"text"`
}

type Identifier struct {
	Type   Type   `json:"type"`
	System string `json:"system"`
	Value  string `json:"value"`
}

type Contained struct {
	ResourceType string       `json:"resourceType"`
	ID           string       `json:"id"`
	Meta         Meta         `json:"meta"`
	Identifiers  []Identifier `json:"identifier"`
}

type ValueCoding struct {
	System  string `json:"system"`
	Code    string `json:"code"`
	Display string `json:"display"`
}

type ValueMoney struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type ValueIdentifier struct {
	System string `json:"system"`
	Value  string `json:"value"`
}

type Extension struct {
	URL             string          `json:"url"`
	ValueCoding     ValueCoding     `json:"valueCoding"`
	ValueMoney      ValueMoney      `json:"valueMoney"`
	ValueIdentifier ValueIdentifier `json:"valueIdentifier"`
}

type Patient struct {
	Reference string `json:"reference"`
}

type Resource struct {
	ResourceType string       `json:"resourceType"`
	ID           string       `json:"id"`
	Meta         Meta         `json:"meta"`
	Containeds   []Contained  `json:"contained"`
	Extensions   []Extension  `json:"extension"`
	Identifiers  []Identifier `json:"identifier"`
	Status       string       `json:"status"`
	Type         Type         `json:"type"`
	SubType      SubType      `json:"subType"`
	Use          string       `json:"use"`
	Patient      Patient      `json:"patient"`
}

type Entry struct {
	Resource Resource
}

type Bundle struct {
	ResourceType string  `json:"resourceType"`
	ID           string  `json:"id"`
	Meta         Meta    `json:"meta"`
	Type         string  `json:"type"`
	Total        int     `json:"total"`
	Links        []Link  `json:"link"`
	Entrys       []Entry `json:"entry"`
}
