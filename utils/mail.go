package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/jung-kurt/gofpdf"
	"github.com/loid-lab/e-commerce-api/models"
	gomail "gopkg.in/mail.v2"
)

func SendMail(data models.EmailData) error {

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("MAIL_FROM"))
	message.SetHeader("To", data.To)
	message.SetHeader("Subject", data.Subject)
	message.SetBody("text/html", data.HTMLBody)

	if data.ImagePath != "" {
		message.Embed(data.ImagePath)
	}

	dialer := gomail.NewDialer(
		os.Getenv("MAIL_HOST"),
		587,
		os.Getenv("MAIL_USER"),
		os.Getenv("MAIL_PASS"),
	)

	return dialer.DialAndSend(message)
}

func CalculateTotalAmount(inv models.Invoice) float64 {
	var total float64
	for _, item := range inv.Items {
		total += float64(item.Quantity) * item.UnitPrice
	}
	return total
}

func GenerateInvoice(invoice models.Invoice) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice #%s", invoice.InvoiceNumber))
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", invoice.Date.Format("2006-01-02")))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", invoice.CustomerName))
	pdf.Ln(12)

	// Table header
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(80, 10, "Description")
	pdf.Cell(30, 10, "Qty")
	pdf.Cell(30, 10, "Unit Price")
	pdf.Cell(30, 10, "Amount")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	for _, item := range invoice.Items {
		pdf.Cell(80, 10, item.Description)
		pdf.Cell(30, 10, fmt.Sprintf("%d", item.Quantity))
		pdf.Cell(30, 10, fmt.Sprintf("$%.2f", item.UnitPrice))
		pdf.Cell(30, 10, fmt.Sprintf("$%.2f", float64(item.Quantity)*item.UnitPrice))
		pdf.Ln(10)
	}

	pdf.Ln(5)
	pdf.Cell(80, 10, "")
	pdf.Cell(30, 10, "Total")
	pdf.Cell(30, 10, fmt.Sprintf("%.2f", invoice.TotalAmount))

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func SendInvoiceEmail(data models.EmailData, cfg models.SMTConfig, invoicePDF []byte) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", data.From)
	msg.SetHeader("To", data.To)
	msg.SetHeader("Subject", data.Subject)
	msg.SetBody("text/html", data.HTMLBody)
	msg.Attach("invoice.pdf", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(invoicePDF)
		return err
	}))

	dialer := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)

	return dialer.DialAndSend(msg)

}
