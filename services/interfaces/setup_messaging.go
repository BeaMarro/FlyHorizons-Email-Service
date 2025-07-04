package interfaces

import "flyhorizons-emailservice/config"

type SetupMessaging interface {
	InitializeRabbitMQ() *config.RabbitMQ
}
