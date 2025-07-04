package main

import (
	"flyhorizons-emailservice/config"
	"flyhorizons-emailservice/internal/health"
	"flyhorizons-emailservice/internal/metrics"
	"flyhorizons-emailservice/services"
	"flyhorizons-emailservice/utils"

	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	healthcfg "github.com/tavsec/gin-healthcheck/config"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	// Initialize RabbitMQ for messaging
	rabbitSetup := &config.SetupMessaging{}
	rabbitMQClient := rabbitSetup.InitializeRabbitMQ()
	config.RabbitMQClient = rabbitMQClient

	defer rabbitMQClient.Connection.Close()
	defer rabbitMQClient.Connection.Channel()

	router := gin.Default()

	// --- Health checks setup ---
	conf := healthcfg.DefaultConfig()
	rabbitMQCheck := health.RabbitMQCheck{}
	healthcheck.New(router, conf, []checks.Check{rabbitMQCheck})

	// --- Metrics setup ---
	metrics.RegisterMetricsRoutes(router, rabbitMQCheck)
	// Utilities
	textUtils := utils.TextUtilities{}

	// Services
	emailSender := &services.EmailSender{}
	emailService := services.NewEmailService(textUtils, rabbitMQClient, emailSender)
	go emailService.StartEmailConsumer()

	// Run the microservice
	router.Run(":8085")
}
