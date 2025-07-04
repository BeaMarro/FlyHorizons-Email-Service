package services_test

import (
	"flyhorizons-emailservice/config"
	"flyhorizons-emailservice/models"
	"flyhorizons-emailservice/models/enums"
	"flyhorizons-emailservice/services"
	"flyhorizons-emailservice/tests/mocks"
	"flyhorizons-emailservice/utils"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestEmailService struct {
}

// Setup
func setupEmailService() (*mocks.MockEmailSender, services.EmailService) {
	mockEmailSender := new(mocks.MockEmailSender)
	mockRabbitMQ := &config.RabbitMQ{
		Connection: nil,
		Channel:    nil,
	}
	textUtils := utils.TextUtilities{}
	emailService := services.NewEmailService(textUtils, mockRabbitMQ, mockEmailSender)
	return mockEmailSender, *emailService
}

func getPassengers() []models.Passenger {
	return []models.Passenger{
		{
			ID:             1,
			FullName:       "John Doe",
			DateOfBirth:    time.Date(1985, 7, 9, 1, 0, 0, 0, time.UTC),
			PassportNumber: "1234",
		},
		{
			ID:             2,
			FullName:       "Jane Doe",
			DateOfBirth:    time.Date(1986, 8, 8, 2, 30, 0, 0, time.UTC),
			PassportNumber: "4321",
		},
	}
}

func getSeats() []models.Seat {
	return []models.Seat{
		{
			Row:       1,
			Column:    "A",
			Available: true,
		},
		{
			Row:       1,
			Column:    "B",
			Available: true,
		},
	}
}

func getLuggageList() []enums.Luggage {
	return []enums.Luggage{enums.SmallBag, enums.Cargo20kg}
}

func getBooking() models.Booking {
	return models.Booking{
		ID:          0,
		UserID:      2,
		FlightCode:  "FR788",
		FlightClass: 1,
		Luggage:     getLuggageList(),
		Seats:       getSeats(),
		Passengers:  getPassengers(),
	}
}

// Service Unit Tests
func TestSendEmailWithBookingSendsEmailToPassengers(t *testing.T) {
	// Arrange
	mockEmailSender, emailService := setupEmailService()
	booking := getBooking()
	mockEmailSender.On("SendEmail",
		mock.Anything,
		[]string{booking.Passengers[0].Email},
		mock.Anything, // Mock the automatically generated email content (HTML)
	).Return(nil)

	// Act
	result := emailService.SendConfirmationEmail(booking)

	// Assert
	mockEmailSender.AssertExpectations(t)
	assert.Nil(t, result) // Nil means no errors occurred
}

func TestSendInvalidEmailWithBookingReturnsError(t *testing.T) {
	// Arrange
	mockEmailSender, emailService := setupEmailService()
	booking := getBooking()
	customError := fmt.Errorf("An error occurred while sending the email.")
	mockEmailSender.On("SendEmail",
		mock.Anything,
		[]string{booking.Passengers[0].Email},
		mock.Anything, // Mock the automatically generated email content (HTML)
	).Return(customError)

	// Act
	result := emailService.SendConfirmationEmail(booking)

	// Assert
	mockEmailSender.AssertExpectations(t)
	assert.Error(t, result)
}
