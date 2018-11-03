package mailer

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go/aws"
)

// Destination represents a single SES bulk email destination.
type Destination struct {
	To   string
	Data interface{}
}

// SESType converts Destination to ses.BulkEmailDestination.
func (d Destination) SESType() ses.BulkEmailDestination {
	j, _ := json.Marshal(d.Data)
	sd := ses.BulkEmailDestination{
		Destination: &ses.Destination{
			ToAddresses: []string{d.To},
		},
		ReplacementTemplateData: aws.String(string(j)),
	}
	return sd
}

// MapToSESType maps array of Destination to ses.BulkEmailDestination.
func MapToSESType(d []Destination) []ses.BulkEmailDestination {
	var ret []ses.BulkEmailDestination
	for _, v := range d {
		ret = append(ret, v.SESType())
	}
	return ret
}
