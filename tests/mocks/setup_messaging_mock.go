package mocks

import (
	"flyhorizons-emailservice/config"
	"flyhorizons-emailservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockSetupMessaging struct {
	mock.Mock
}

var _ interfaces.SetupMessaging = (*MockSetupMessaging)(nil)

func (m *MockSetupMessaging) InitializeRabbitMQ() *config.RabbitMQ {
	args := m.Called()
	return args.Get(0).(*config.RabbitMQ)
}
