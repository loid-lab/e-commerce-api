package models

type EmailData struct {
	To        string
	Subject   string
	HTMLBody  string
	ImagePath string // optional
}
