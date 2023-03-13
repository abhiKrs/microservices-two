package email

import (
	"bytes"
	"html/template"

	constants "notification/app/constants"
)

type TemplateBody struct {
	Subject          string
	PreheaderMessage string
	Heading          string
	TextBody         string
	Name             string
	URL              string
	Matter           string
	ButtonAction     string
	Footnote         string
}

func (tb *TemplateBody) MagicLinkTemplate(linkType constants.DBMagicLinkType, redirectLink string) {
	if linkType == constants.Register {
		tb.Subject = "Welcome to Logfire"
		tb.PreheaderMessage = ":" //"Registration Process-"
		tb.TextBody = "Click the button below to securely complete your registration process with LogFire."
		tb.Heading = "Logfire Invite Magic-Link."
		tb.Matter = "We at LogFire invite you to the world of dynamic logs."
		tb.Name = "there"
		tb.ButtonAction = "Verify Registration"
		tb.URL = redirectLink
		tb.Footnote = "This link is provided for sigle time use only. The link will expire in 1 hour."
	} else if linkType == constants.Login {
		tb.Subject = "Login to Logfire"
		tb.PreheaderMessage = ":" //"Login Magic-Link Inside-"
		tb.TextBody = "Click the button below to securely complete your login process with LogFire."
		tb.Heading = "Logfire Login Magic-Link."
		tb.Matter = "Please continue to the world of dynamic logs."
		tb.Name = "there"
		tb.ButtonAction = "Login"
		tb.URL = redirectLink
		tb.Footnote = "This link is provided for sigle time use only. The link will expire in 1 hour."
	} else if linkType == constants.ResetPassword {
		tb.Subject = "Reset Password to Logfire"
		tb.PreheaderMessage = ":" //"Password-Reset Request-"
		tb.TextBody = "Click the button below to securely complete your password-reset process with LogFire."
		tb.Heading = "Logfire Reset Password Magic-Link."
		tb.Matter = "Please continue to reset your password."
		tb.Name = "there"
		tb.ButtonAction = "Reset Password"
		tb.URL = redirectLink
		tb.Footnote = "This link is provided for sigle time use only. The link will expire in 1 hour."
	}
}

func (tb *TemplateBody) SuccessTemplate(successType constants.SuccessType, redirectLink string) {
	if successType == constants.Onboarding {
		tb.Subject = "Successful Onboarding to Logfire"
		tb.PreheaderMessage = ":" //"Password-Reset Request-"
		tb.TextBody = "You have successfuly completed your onboarding with Logfire. We have added you to our Access-Request list./nWe will contact you soon as we launch."
		tb.Heading = "Logfire Welcome Message."
		tb.Matter = "We at LogFire invite you to the world of dynamic logs."
		tb.Name = "there"
		tb.ButtonAction = "Visit Logfire"
		tb.URL = redirectLink
		tb.Footnote = "........"
	} else if successType == constants.PasswordReset {
		tb.Subject = "Successful Password-reset for Logfire"
		tb.PreheaderMessage = ":" //"Password-Reset Request-"
		tb.TextBody = "You have successfuly reset your password. Please continue to login with your new password."
		tb.Heading = "Logfire Success Message."
		tb.Matter = "You have successfully reset your password."
		tb.Name = "there"
		tb.ButtonAction = "Visit Logfire"
		tb.URL = redirectLink
		tb.Footnote = "........"
	} else if successType == constants.PasswordSet {
		tb.Subject = "Successful Password-Set for Logfire"
		tb.PreheaderMessage = ":" //"Password-Reset Request-"
		tb.TextBody = "You have successfuly set your password. Please continue to login with your new password."
		tb.Heading = "Logfire Success Message."
		tb.Matter = "You have successfully set your password."
		tb.Name = "there"
		tb.ButtonAction = "Visit Logfire"
		tb.URL = redirectLink
		tb.Footnote = "........"
	}
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)

	if err != nil {
		return "nil", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()

	return body, nil

}
