package utils

import (
	"fmt"

	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
)

func GenerateSendInvoice(inv models.Invoice, emailData models.EmailData) error {
	pdf, err := GenerateInvoice(inv)
	if err != nil {
		return fmt.Errorf("failed to generate invoice: %v", err)
	}

	smtpConfig := models.SMTConfig{
		SMTPHost: initializers.Env.SMTPHost,
		SMTPPort: initializers.Env.SMTPPort,
		SMTPUser: initializers.Env.SMTPUser,
		SMTPPass: initializers.Env.SMTPPass,
	}

	err = SendInvoiceEmail(emailData, smtpConfig, pdf)
	if err != nil {
		return fmt.Errorf("failed to send invoice email: %v", err)
	}

	return nil
}
