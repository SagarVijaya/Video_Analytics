package mailsender

import (
	"fmt"
	"log"
	"net/smtp"
	"videoanalytics/apps/models"
	"videoanalytics/config"
	// adjust path as needed
)

func SendFailureAlert(pClickInfo models.AdClickRate, pdbErr error) {
	subject := "DB INSERT Failed - Redis Fallback Triggered"

	body := fmt.Sprintf(`
Database insert failed for AdClick:
--------------------------------------------------
Ad ID           : %d
Clicked At (UTC): %s
IP Address      : %s
Playback Second : %f
--------------------------------------------------

Error message:
%s
Please investigate the database container or network connectivity.
`, pClickInfo.AdID, pClickInfo.ClickedAt.UTC(), pClickInfo.IP, pClickInfo.PlaybackSecond, pdbErr.Error())

	if err := SendAlertMail(subject, body); err != nil {
		log.Println("Mailer Error:", err)
	}
}

// SendAlertMail sends an email with the given subject and body to the specified recipient.
func SendAlertMail(subject, body string) error {
	from := config.GetConfig().Mail.Host
	pass := config.GetConfig().Mail.Pass
	smtpHost := config.GetConfig().Mail.Host
	smtpPort := config.GetConfig().Mail.Port

	// Compose the message
	msg := "From: " + from + "\n" +
		"To: " + config.GetConfig().Mail.To + "\n" +
		"Subject: " + subject + "\n" +
		body

	// Set up authentication info
	auth := smtp.PlainAuth("", from, pass, smtpHost)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{config.GetConfig().Mail.To}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send alert mail: %w", err)
	}
	return nil
}
