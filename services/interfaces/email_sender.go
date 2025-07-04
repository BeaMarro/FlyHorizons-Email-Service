package interfaces

type EmailSender interface {
	SendEmail(from string, to []string, message []byte) error
}
