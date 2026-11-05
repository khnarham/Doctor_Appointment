package mailer

import (
	"context"
	"fmt"
	"os"
	"strings"

	brevo "github.com/getbrevo/brevo-go/lib"
)

type brevoMailer struct {
	client *brevo.APIClient
	from   string
}

func NewBrevoMailer(apiKey string, from string) *brevoMailer {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)

	return &brevoMailer{
		client: brevo.NewAPIClient(cfg),
		from:   from,
	}
}

func (b *brevoMailer) SendVerificationToken(
	ctx context.Context,
	to string,
	token string,
) error {


	     
	frontendUri := os.Getenv("FRONTEND_URI")
	fmt.Println("DEBUG FRONTEND_URI FROM RAILWAY:", os.Getenv("FRONTEND_URI"))
    if frontendUri == "" {
    frontendUri = "http://localhost:3000"
    }
    frontendUri = strings.TrimRight(frontendUri, "/")


verificationURL := fmt.Sprintf(
    "%s/verify-email?token=%s",
    frontendUri,
    token,
)
	html := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
            <body>
                <h2>Email Verification</h2>

                <p>
                    Please click the link below to verify your email.
                </p>

                <a href="%s">
                    Verify Email
                </a>
            </body>
        </html>
    `, verificationURL)

	params := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Email: b.from,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: to,
			},
		},
		Subject:     "Verify Your Email",
		HtmlContent: html,
	}

	_, _, err := b.client.TransactionalEmailsApi.SendTransacEmail(ctx, params)
	if err != nil {
		return fmt.Errorf("brevo send failed: %w", err)
	}

	return nil
}