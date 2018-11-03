package mailer

import "encoding/json"

// Template represents a SES template.
type Template struct {
	Name        string
	Subject     string
	Text        string
	HTML        string
	DefaultData interface{}
}

// DefaultDataJSON returns the stringify json of the default data.
func (t Template) DefaultDataJSON() string {
	d, _ := json.Marshal(t.DefaultData)
	return string(d)
}
