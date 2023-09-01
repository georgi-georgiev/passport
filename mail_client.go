package passport

import (
	"context"
	"io"
	"net/http"

	brevo "github.com/getbrevo/brevo-go/lib"
	"go.uber.org/zap"
)

type MailClient struct {
	senderEmail string
	br          *brevo.APIClient
	log         *zap.Logger
}

func NewMailCleint(conf *Config, log *zap.Logger) *MailClient {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", conf.Mail.ApiKey)
	br := brevo.NewAPIClient(cfg)
	return &MailClient{br: br, senderEmail: conf.Mail.SenderEmail, log: log}
}

func (mc *MailClient) Send(ctx context.Context, subject string, body string, email string) error {
	_, resp, err := mc.br.TransactionalEmailsApi.SendTransacEmail(ctx, brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Email: mc.senderEmail,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: email,
			},
		},
		Subject:     subject,
		TextContent: body,
	})

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {

		resp.Body = io.NopCloser(ReusableReader((resp.Body)))
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		mc.log.Info("Email not sent successfully")
	}

	return nil
}
