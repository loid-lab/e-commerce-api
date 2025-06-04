package models

type EmailData struct {
	To        string
	Subject   string
	HTMLBody  string
	ImagePath string
	From      string
	SMTConfig SMTConfig
}

type SMTConfig struct {
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
}
