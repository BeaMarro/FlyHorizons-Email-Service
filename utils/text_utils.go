package utils

type TextUtilities struct {
}

func JoinStrings(strings []string, separator string) string {
	result := ""
	for index, string := range strings {
		if index > 0 {
			result += separator
		}
		result += string
	}
	return result
}
