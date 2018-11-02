// Package mailer reference:
// https://godoc.org/github.com/aws/aws-sdk-go-v2/service/ses
package mailer

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

// MaxReceivers specifies the maximum number of receivers per bulk email.
const MaxReceivers int = 50

// Emailer represents the ses client for creating template and sending bulk email.
type Emailer struct {
	From   string
	Client *ses.SES
}

// New creates new Emailer instance.
func New(from string) (*Emailer, error) {
	// load configs
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}
	cfg.Region = endpoints.UsEast1RegionID

	// init ses clients
	return &Emailer{
		From:   from,
		Client: ses.New(cfg),
	}, nil
}

// CreateTemplate creates an email template. Existing template of the same name
// will be over-written.
func (e *Emailer) CreateTemplate(template Template) (*ses.CreateTemplateOutput, error) {
	fmt.Println("Creating template...")
	t := ses.Template{
		TemplateName: aws.String(template.Name),
		SubjectPart:  aws.String(template.Subject),
		TextPart:     aws.String(template.Text),
		HtmlPart:     aws.String(template.HTML),
	}
	input := ses.CreateTemplateInput{Template: &t}
	req := e.Client.CreateTemplateRequest(&input)
	res, err := req.Send()
	if err != nil {
		if strings.HasPrefix(err.Error(), ses.ErrCodeAlreadyExistsException) {
			_, err = e.DeleteTemplate(template.Name)
			if err != nil {
				return nil, err
			}
			fmt.Println("Resend create template request...")
			req2 := e.Client.CreateTemplateRequest(&input)
			res, err = req2.Send()
			if err != nil {
				return nil, err
			}
			return res, nil
		}
		return nil, err
	}
	return res, nil
}

// DeleteTemplate deletes the template of the given name.
func (e *Emailer) DeleteTemplate(temaplateName string) (*ses.DeleteTemplateOutput, error) {
	fmt.Println("Deleting template...")
	t := ses.DeleteTemplateInput{
		TemplateName: aws.String(temaplateName),
	}
	req := e.Client.DeleteTemplateRequest(&t)
	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Send sends bulk emails in batch.
func (e *Emailer) Send(template Template, destinations []ses.BulkEmailDestination) ([]*ses.SendBulkTemplatedEmailOutput, error) {
	total := len(destinations)
	var allRes []*ses.SendBulkTemplatedEmailOutput
	for i := 0; i < total; i += MaxReceivers {
		end := i + MaxReceivers
		if end > total {
			end = total
		}
		des := destinations[i:end]
		res, err := e.SendSingle(template, des)
		if err != nil {
			return nil, err
		}
		allRes = append(allRes, res)
		if end < total {
			time.Sleep(3 * time.Second)
		}
	}
	return allRes, nil
}

// SendSingle sends a single bulk email to every address in destinations.
func (e *Emailer) SendSingle(template Template, destinations []ses.BulkEmailDestination) (*ses.SendBulkTemplatedEmailOutput, error) {
	fmt.Println("Sending bulk email...")
	input := &ses.SendBulkTemplatedEmailInput{
		Source:               aws.String(e.From),
		ConfigurationSetName: aws.String("log"),
		Template:             aws.String(template.Name),
		Destinations:         destinations,
		DefaultTemplateData:  aws.String(template.DefaultData),
	}

	req := e.Client.SendBulkTemplatedEmailRequest(input)
	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// BuildDest helps to build a destination target.
func BuildDest(to, templateData string) ses.BulkEmailDestination {
	d := ses.BulkEmailDestination{
		Destination: &ses.Destination{
			ToAddresses: []string{to},
		},
		ReplacementTemplateData: aws.String(templateData),
	}
	return d
}
