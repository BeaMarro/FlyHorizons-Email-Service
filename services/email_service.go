package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mime/quotedprintable"
	"os"

	"flyhorizons-emailservice/config"
	"flyhorizons-emailservice/models"
	"flyhorizons-emailservice/services/interfaces"
	"flyhorizons-emailservice/utils"

	"github.com/joho/godotenv"
)

type EmailService struct {
	textUtilities  utils.TextUtilities
	rabbitMQClient *config.RabbitMQ
	emailSender    interfaces.EmailSender
}

func NewEmailService(utils utils.TextUtilities, client *config.RabbitMQ, sender interfaces.EmailSender) *EmailService {
	return &EmailService{
		textUtilities:  utils,
		rabbitMQClient: client,
		emailSender:    sender,
	}
}

func (emailService *EmailService) StartEmailConsumer() {
	channel := emailService.rabbitMQClient.Channel

	messages, err := channel.Consume(
		"booking.confirmed", // Queue name
		"",                  // Consumer tag
		true,                // Auto acknowledge
		false,               // Exclusive
		false,               // No-local
		false,               // No-wait
		nil,                 // Arguments
	)
	if err != nil {
		log.Fatalf("An error occurred while registering the consumer: %v", err)
	}

	log.Printf("Started consumer for RabbitMQ queue: %s", "booking.confirmed")
	log.Printf("Channel: RabbitMQ channel is active")

	go func() {
		for message := range messages {
			log.Printf("Messages:")
			log.Printf("Body: %s", string(message.Body))

			var booking models.Booking
			if err := json.Unmarshal(message.Body, &booking); err != nil {
				log.Printf("An error occurred while converting the JSON to a Booking object: %v", err)
				continue
			}

			emailService.SendConfirmationEmail(booking)
		}
		log.Printf("Consumer for queue %s stopped", "booking.confirmed")
	}()

	log.Println("Waiting for a Booking to be created. To exit press CTRL+C")
	forever := make(chan bool)
	<-forever
}

// createMIMEEmail builds the full multipart email with HTML body and inline image (qr code)
func (emailService *EmailService) createMIMEEmail(from, to, subject, htmlBody string, qrPNG []byte) ([]byte, error) {
	var buf bytes.Buffer

	boundary := "MY-MULTIPART-BOUNDARY"

	// Headers
	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/related; boundary=%s\r\n", boundary))
	buf.WriteString("\r\n")

	// HTML part
	buf.WriteString("--" + boundary + "\r\n")
	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	buf.WriteString("\r\n")

	qpWriter := quotedprintable.NewWriter(&buf)
	if _, err := qpWriter.Write([]byte(htmlBody)); err != nil {
		return nil, err
	}
	qpWriter.Close()
	buf.WriteString("\r\n")

	// Image part
	buf.WriteString("--" + boundary + "\r\n")
	buf.WriteString("Content-Type: image/png\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString("Content-ID: <qr-code>\r\n")
	buf.WriteString("Content-Disposition: inline; filename=\"qrcode.png\"\r\n")
	buf.WriteString("\r\n")

	encoded := base64.StdEncoding.EncodeToString(qrPNG)
	// Break base64 lines to max 76 chars per RFC
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		buf.WriteString(encoded[i:end] + "\r\n")
	}

	// Closing boundary
	buf.WriteString("--" + boundary + "--\r\n")

	return buf.Bytes(), nil
}

func (emailService *EmailService) SendConfirmationEmail(booking models.Booking) error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	from := os.Getenv("SENDER_EMAIL")
	to := booking.Passengers[0].Email

	var passengerInfo []string
	for _, passenger := range booking.Passengers {
		passengerInfo = append(passengerInfo, passenger.String())
	}
	passengerString := utils.JoinStrings(passengerInfo, ", ")

	var seatInfo []string
	for _, seat := range booking.Seats {
		seatInfo = append(seatInfo, seat.String())
	}
	seatString := utils.JoinStrings(seatInfo, ", ")

	var luggageInfo []string
	for _, luggage := range booking.Luggage {
		luggageInfo = append(luggageInfo, string(luggage))
	}
	luggageString := utils.JoinStrings(luggageInfo, ", ")

	qrContent := fmt.Sprintf("BookingID: %d, Flight: %s", booking.ID, booking.FlightCode)

	// Generate QR code PNG bytes
	qrPNG, err := utils.GenerateQRCodePNG(qrContent)
	if err != nil {
		return fmt.Errorf("Failed to generate QR code PNG: %v", err)
	}

	// HTML references the qr-code by cid
	subject := fmt.Sprintf("Booking Confirmation for FlyHorizons Flight %s", booking.FlightCode)
	body := fmt.Sprintf(`
	<html>
	<body>
		<p>Dear %s,</p>
		<p>✈️ Thank you for making your booking with FlyHorizons.</p>
		<p>Your booking has been successfully confirmed for flight <b>%s</b>.</p>
		<p>--------------------------------------------------------------------------------------------</p>
		<p><b>Please find your booking information below:</b></p>
		<p><b>Passengers:</b> %s</p>
		<p><b>Seats:</b> %s</p>
		<p><b>Flight Class:</b> %s</p>
		<p><b>Luggage:</b> %s</p>
		<p>--------------------------------------------------------------------------------------------</p>
		<p><b>Scan this QR code at the gate:</b><br/><img src="cid:qr-code" alt="QR Code"/></p>
		<p>--------------------------------------------------------------------------------------------</p>
		<p>We wish you a nice flight!</p>
		<p>Best regards,</p>
		<p>FlyHorizons European Airlines</p>
	</body>
	</html>`, booking.Passengers[0].FullName, booking.FlightCode, passengerString, seatString, booking.FlightClass.String(), luggageString)

	// Create MIME email with embedded image
	message, err := emailService.createMIMEEmail(from, to, subject, body, qrPNG)
	if err != nil {
		return fmt.Errorf("Failed to create MIME email: %v", err)
	}

	err = emailService.emailSender.SendEmail(from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("Failed to send email to %s: %v", to, err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
