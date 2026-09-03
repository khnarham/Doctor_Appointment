package Mailer

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

type resendMailer struct {
	client *resend.Client
	from   string
}

func NewResendMailer(apiKey string, from string) *resendMailer {
	client := resend.NewClient(apiKey)

	return &resendMailer{
		client: client,
		from:   from,
	}
}

func (r *resendMailer) SendVerificationToken(
	ctx context.Context,
	to string,
	token string,
) error {

	verificationURL := fmt.Sprintf(
		"http://localhost:3000/verify-email?token=%s",
		token,
	)

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<body>
				<h2>Email Verification</h2>

				<p>
					Please click the button below to verify your email.
				</p>

				<a href="%s">
					Verify Email
				</a>
			</body>
		</html>
	`, verificationURL)

	params := &resend.SendEmailRequest{
		From:    r.from,
		To:      []string{to},
		Subject: "Verify Your Email",
		Html:    html,
	}

	_, err := r.client.Emails.Send(params)

	if err != nil {
		return fmt.Errorf("resend send failed: %w", err)
	}

	return nil
}