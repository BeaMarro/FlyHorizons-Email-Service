package mocks

import (
	"flyhorizons-emailservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockEmailSender struct {
	mock.Mock
}

var _ interfaces.EmailSender = (*MockEmailSender)(nil)

func (m *MockEmailSender) SendEmail(from string, to []string, message []byte) error {
	args := m.Called(from, to, message)
	return args.Error(0)
}
