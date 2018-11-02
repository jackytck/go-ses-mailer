package mailer

// Template represents a SES template.
type Template struct {
	Name        string
	Subject     string
	Text        string
	HTML        string
	DefaultData string
}
