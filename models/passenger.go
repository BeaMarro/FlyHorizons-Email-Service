package models

import (
	"fmt"
	"time"
)

type Passenger struct {
	ID             int       `json:"id"`
	FullName       string    `json:"full_name"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	PassportNumber string    `json:"passport_number"`
	Email          string    `json:"email"`
}

func (passenger Passenger) String() string {
	return fmt.Sprintf("Full Name: %s, Email: %s, Passport Number: %s", passenger.FullName, passenger.Email, passenger.PassportNumber)
}
