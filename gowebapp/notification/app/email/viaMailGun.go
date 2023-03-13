package email

import (
	"context"
	"time"
	// "fmt"

	"notification/app/config"
	log "notification/app/utility"
	constants "notification/app/constants"

	"github.com/mailgun/mailgun-go/v4"
)

func SendEmailViaMailGun(email string, token *string, messageType constants.EmailType, linkType *constants.DBMagicLinkType, successType *constants.SuccessType, cfg *config.Config) error {

	// Your available domain names can be found here:
	// (https://app.mailgun.com/app/domains)
	var yourDomain string = cfg.Api.MAILGUN_DOMAIN // e.g. mg.yourcompany.com
	log.DebugLogger.Println(yourDomain)

	// You can find the Private API Key in your Account Menu, under "Settings":
	// (https://app.mailgun.com/app/account/security)
	var privateAPIKey string = cfg.Api.MAILGUN_KEY
	log.DebugLogger.Println(privateAPIKey)
	// Create an instance of the Mailgun Client

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
	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	sender := cfg.Api.MAILER_EMAIL
	subject := templateData.Subject
	recipient := email

	message := mg.NewMessage(sender, subject, "", recipient)

	message.SetHtml(HtmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 20)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}

	log.InfoLogger.Printf("ID: %s Resp: %s\n", id, resp)
	return nil

}
