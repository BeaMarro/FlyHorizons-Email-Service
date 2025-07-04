package utils_test

import (
	"flyhorizons-emailservice/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TextUtilsTest struct {
}

func TestJoinValidStringsReturnsJoinedString(t *testing.T) {
	// Arrange
	string1 := "Hello"
	string2 := "World!"
	strings := []string{string1, string2}
	delimeter := ", "

	// Act
	joinedString := utils.JoinStrings(strings, delimeter)

	// Assert
	assert.Equal(t, "Hello, World!", joinedString)
}
