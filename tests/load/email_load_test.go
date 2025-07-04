package service_test

import (
	"flyhorizons-emailservice/models"
	"flyhorizons-emailservice/models/enums"
	"flyhorizons-emailservice/services"
	load_test_utils "flyhorizons-emailservice/tests/load/utils"
	"flyhorizons-emailservice/utils"
	"net/http"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// TODO: I think it is sufficient to leave the load testing for both on the BookingService
type EmailLoadTest struct {
	loadTestUtils load_test_utils.LoadTestUtils
}

// Fake booking for testing
func getFakeBooking() models.Booking {
	return models.Booking{
		ID:          1,
		UserID:      123,
		FlightCode:  "FR100",
		FlightClass: enums.Business,
		Luggage:     []enums.Luggage{enums.CabinBag, enums.Cargo20kg},
		Seats: []models.Seat{
			{Row: 2, Column: "B", Available: true},
		},
		Passengers: []models.Passenger{
			{
				ID:             1,
				FullName:       "Alice Smith",
				Email:          "jamics.mail@gmail.com",
				DateOfBirth:    time.Date(1990, 3, 15, 0, 0, 0, 0, time.UTC),
				PassportNumber: "5678",
			},
		},
	}
}

// TODO: Fix this to use the RabbitMQ and somehow mock this
func TestLoad_SendRealEmail(t *testing.T) {
	// Initialize utilities and services
	textUtils := utils.TextUtilities{}
	emailSender := &services.EmailSender{}
	emailService := services.NewEmailService(textUtils, nil, emailSender)

	// Get fake booking
	booking := getFakeBooking()

	// Initialize metrics collector
	loadTest := EmailLoadTest{loadTestUtils: load_test_utils.LoadTestUtils{}}
	var metrics vegeta.Metrics

	// Send email and measure latency
	start := time.Now()
	err := emailService.SendConfirmationEmail(booking)
	latency := time.Since(start)

	// Create result for metrics
	result := vegeta.Result{
		Code:      http.StatusOK,
		Timestamp: time.Now(),
		Attack:    "Direct Email Test",
		Latency:   latency,
	}

	if err != nil {
		result.Code = http.StatusInternalServerError
		result.Error = err.Error()
	}

	// Store and evaluate results
	metrics.Add(&result)
	metrics.Close()

	loadTest.loadTestUtils.LogMetrics(t, &metrics)
	loadTest.loadTestUtils.EvaluateMetricsSuccess(t, &metrics)
	loadTest.loadTestUtils.SaveToJSON(t, &metrics)
}
