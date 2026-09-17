package mailer

import (
	"context"
	"fmt"

	brevo "github.com/getbrevo/brevo-go/lib"
)

type BrevoReceptionistMailer struct {
	client  *brevo.APIClient
	from    string
	baseURL string
}

func NewBrevoReceptionistMailer(apiKey string, from string, baseURL string) *BrevoReceptionistMailer {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)

	// Fallback to default domain if baseURL is not provided
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	return &BrevoReceptionistMailer{
		client:  brevo.NewAPIClient(cfg),
		from:    from,
		baseURL: baseURL,
	}
}

// SendDoctorInviteEmail sends invitation link to candidate receptionist
func (b *BrevoReceptionistMailer) SendDoctorInviteEmail(
	ctx context.Context,
	toEmail string,
	token string,
) error {

	// Direct invitation route configured for receptionist registration
	inviteURL := fmt.Sprintf(
		"%s/accept-invite/receptionist-register?token=%s",
		b.baseURL,
		token,
	)

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<style>
					body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
					.container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px; }
					.button { display: inline-block; padding: 12px 24px; background-color: #007bff; color: #ffffff; text-decoration: none; border-radius: 5px; font-weight: bold; margin-top: 15px; }
					.footer { margin-top: 25px; font-size: 12px; color: #777; }
				</style>
			</head>
			<body>
				<div class="container">
					<h2>Doctor Receptionist Invitation</h2>
					<p>You have been invited to join as a Receptionist in our Clinic Appointment System.</p>
					<p>Please click the button below to accept the invitation and complete your account registration:</p>
					
					<a href="%s" class="button" target="_blank">
						Accept Invitation & Register
					</a>

					<p style="margin-top: 20px; font-size: 13px;">
						If the button doesn't work, copy and paste this URL into your browser: <br>
						<a href="%s">%s</a>
					</p>

					<div class="footer">
						<p>If you were not expecting this invitation, you can safely ignore this email.</p>
					</div>
				</div>
			</body>
		</html>
	`, inviteURL, inviteURL, inviteURL)

	params := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  "Appointment System Support",
			Email: b.from,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: toEmail,
			},
		},
		Subject:     "Invitation to Join as Receptionist",
		HtmlContent: html,
	}

	_, _, err := b.client.TransactionalEmailsApi.SendTransacEmail(ctx, params)
	if err != nil {
		return fmt.Errorf("brevo receptionist email dispatch failed: %w", err)
	}

	return nil
}