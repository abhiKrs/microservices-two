package email

import (
	"fmt"
	log "notification/app/utility"

	"notification/app/config"
	constants "notification/app/constants"

	// "github.com/go-mail/mail"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func SendEmailviaAws(email string, token *string, messageType constants.EmailType, linkType *constants.DBMagicLinkType, successType *constants.SuccessType, cfg *config.Config) error {
	var Sender = cfg.Api.MAILER_EMAIL
	// log.Println(cfg.Api.MAILER_EMAIL)
	var templateData = TemplateBody{}
	if messageType == constants.MagicLinkMessage {

		LinkUrl := cfg.Client.BASE_URL
		log.DebugLogger.Println((LinkUrl))
		redirectLink := "https://" + LinkUrl + "/link/?token=" + *token

		// var TextBody = "This email was sent from Logfire. Please click on the link below to complete your registeration."

		templateData.MagicLinkTemplate(*linkType, redirectLink)
		log.DebugLogger.Println(templateData)
	} else if messageType == constants.SuccessMessage {

		LinkUrl := cfg.Client.BASE_URL
		log.DebugLogger.Println((LinkUrl))
		redirectLink := "https://" + LinkUrl

		// var TextBody = "This email was sent from Logfire. Please click on the link below to complete your registeration."

		templateData.SuccessTemplate(*successType, redirectLink)
		log.DebugLogger.Println(templateData)
	}

	HtmlBody, err := ParseTemplate("./public/emailTemplate/base.html", templateData)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	// ---------------------------------

	var CharSet = "UTF-8"

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(cfg.Api.AWS_ACCESS_KEY_ID, cfg.Api.AWS_SECRET_ACCESS_KEY, ""),
	},
	)
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	// Create an SES session.
	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(templateData.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(templateData.Subject),
			},
		},
		Source: aws.String(Sender),
	}

	result, err := svc.SendEmail(input)
	if err != nil {
		// log.Println
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.ErrorLogger.Println(ses.ErrCodeMessageRejected, aerr.Error())
				// return
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.ErrorLogger.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.ErrorLogger.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.ErrorLogger.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.ErrorLogger.Println(err.Error())
		}
		log.ErrorLogger.Println(err)
		return err
	}

	fmt.Println("Email Sent to address: " + email)
	fmt.Println(result)
	return nil
}

// func formatMail() {
// 	// auth = smtp.PlainAuth("", "dhanush@geektrust.in", "password", "smtp.gmail.com")
// 	templateData := struct {
// 		Name string
// 		URL  string
// 	}{
// 		Name: "Dhanush",
// 		URL:  "http://geektrust.in",
// 	}
// 	r := NewRequest([]string{"junk@junk.com"}, "Hello Junk!", "Hello, World!")
// 	// err := r.ParseTemplate("template.html", templateData)
// 	if err := r.ParseTemplate("template.html", templateData); err == nil {
// 		ok, _ := r.SendEmail()
// 		fmt.Println(ok)
// 	}

// }

// //Request struct
// type Request struct {
// 	from    string
// 	to      []string
// 	subject string
// 	body    string
// }

// func NewRequest(to []string, subject, body string) *Request {
// 	return &Request{
// 		to:      to,
// 		subject: subject,
// 		body:    body,
// 	}
// }

// func (r *Request) SendEmail() (bool, error) {
// 	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
// 	subject := "Subject: " + r.subject + "!\n"
// 	msg := []byte(subject + mime + "\n" + r.body)
// 	addr := "smtp.gmail.com:587"

// 	if err := smtp.SendMail(addr, auth, "dhanush@geektrust.in", r.to, msg); err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
// 	t, err := template.ParseFiles(templateFileName)
// 	if err != nil {
// 		return err
// 	}
// 	buf := new(bytes.Buffer)
// 	if err = t.Execute(buf, data); err != nil {
// 		return err
// 	}
// 	r.body = buf.String()
// 	return nil
// }

// -----------------------------------------------------------------------------------
// type Recipient struct {
// 	toEmails  []string
// 	ccEmails  []string
// 	bccEmails []string
// }

// func sendEmailviaAwsNew(email *string, token string, cfg *config.Config) error {
// 	var LinkUrl string
// 	if cfg.Api.ENV == "dev" {
// 		LinkUrl = cfg.Client.BASE_URL
// 	} else {
// 		LinkUrl = cfg.Client.HOST + ":" + cfg.Client.PORT
// 	}

// 	msg := "https://" + LinkUrl + "/?email=" + *email + "&token=" + token

// 	message := "<h1>LogFire Invite Email </h1><p>This email was sent from " +
// 		"<a href='https://logfire.sh'>Logfire.</a> Please click on the link below to complete your registeration.</p> " +
// 		"<p>" + msg + "</p>"
// 	subject := "Welcome to LogFire"
// 	fromEmail := cfg.Api.MAILER_EMAIL

// 	recipient := Recipient{
// 		toEmails:  []string{*email},
// 		ccEmails:  []string{},
// 		bccEmails: []string{},
// 	}
// 	log.Println(recipient.toEmails)

// 	err := SendEmailSES(message, subject, fromEmail, recipient, cfg)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func SendEmailSES(messageBody string, subject string, fromEmail string, recipient Recipient, cfg *config.Config) error {

// 	// create new AWS session
// 	sess, err := session.NewSession(&aws.Config{
// 		Region:      aws.String("us-east-1"),
// 		Credentials: credentials.NewStaticCredentials(cfg.Api.AWS_ACCESS_KEY_ID, cfg.Api.AWS_SECRET_ACCESS_KEY, ""),
// 	})

// 	if err != nil {
// 		log.Println("Error occurred while creating aws session", err)
// 		return err
// 	}
// 	// set to section
// 	var recipients []*string
// 	for _, r := range recipient.toEmails {
// 		recipient := r
// 		recipients = append(recipients, &recipient)
// 	}

// 	// set cc section
// 	var ccRecipients []*string
// 	if len(recipient.ccEmails) > 0 {
// 		for _, r := range recipient.ccEmails {
// 			ccrecipient := r
// 			ccRecipients = append(ccRecipients, &ccrecipient)
// 		}
// 	}

// 	// set bcc section
// 	var bccRecipients []*string
// 	if len(recipient.bccEmails) > 0 {
// 		for _, r := range recipient.bccEmails {
// 			bccrecipient := r
// 			recipients = append(recipients, &bccrecipient)
// 		}
// 	}
// 	log.Println(recipients)
// 	// create an SES session.
// 	svc := ses.New(sess)

// 	// Assemble the email.
// 	input := &ses.SendEmailInput{

// 		// Set destination emails
// 		Destination: &ses.Destination{
// 			CcAddresses:  ccRecipients,
// 			ToAddresses:  recipients,
// 			BccAddresses: bccRecipients,
// 		},

// 		// Set email message and subject
// 		Message: &ses.Message{
// 			Body: &ses.Body{
// 				Html: &ses.Content{
// 					Charset: aws.String("UTF-8"),
// 					Data:    aws.String(messageBody),
// 				},
// 			},

// 			Subject: &ses.Content{
// 				Charset: aws.String("UTF-8"),
// 				Data:    aws.String(subject),
// 			},
// 		},

// 		// send from email
// 		Source: aws.String(fromEmail),
// 	}

// 	// Call AWS send email function which internally calls to SES API
// 	_, err = svc.SendEmail(input)
// 	if err != nil {
// 		log.Println("Error sending mail - ", err)
// 		return err
// 	}

// 	log.Println("Email sent successfully to: ", recipient.toEmails)
// 	return nil
// }

// ---------------------------------------------------------------------------------------
// func sendMagicLinkviaGoogle(email string, token string, cfg *config.Config) error {
// 	var LinkUrl string
// 	if cfg.Api.ENV == "dev" {
// 		LinkUrl = cfg.Client.BASE_URL
// 	} else {
// 		LinkUrl = cfg.Client.HOST + ":" + cfg.Client.PORT
// 	}
// 	m := mail.NewMessage()
// 	m.SetHeader("From", cfg.Api.MAILER_EMAIL)

// 	m.SetHeader("To", email)

// 	// m.SetAddressHeader("Cc", "oliver.doe@example.com", "Oliver")

// 	m.SetHeader("Subject", "Hello!")

// 	msg := "https://" + LinkUrl + "/?email=" + email + "&token=" + token
// 	log.Println("-------inside sendMagicLink--------")
// 	m.SetBody("text/html", fmt.Sprintf("Hello <b>User</b> and click the link below to register: <p>%v</p>", msg))

// 	// m.Attach("lolcat.jpg")

// 	d := mail.NewDialer(
// 		cfg.Api.MAILER_SMTP_SERVER,
// 		25,
// 		cfg.Api.MAILER_USERNAME,
// 		cfg.Api.MAILER_PASSWORD,
// 	)

// 	// Send the email to Kate, Noah and Oliver.
// 	log.Println("inside sendMagicLink")
// 	if err := d.DialAndSend(m); err != nil {
// 		// log.Println(err.Error())
// 		return err
// 		// panic(err)

// 	}

// 	return nil

// }
