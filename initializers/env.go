package initializers

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/loid-lab/e-commerce-api/models"
)

var Env models.SMTConfig

func LoadEnv() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load Environment variable %s", err)
	}

	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: must be an integer. Got %s", smtpPortStr)
	}

	Env = models.SMTConfig{
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: smtpPort,
		SMTPUser: os.Getenv("SMTP_USER"),
		SMTPPass: os.Getenv("SMTP_PASS"),
	}

}
