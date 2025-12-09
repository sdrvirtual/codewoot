package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var brValidDDDs = map[int]struct{}{
	11: {}, 12: {}, 13: {}, 14: {}, 15: {}, 16: {}, 17: {}, 18: {}, 19: {},
	21: {}, 22: {}, 24: {}, 27: {}, 28: {},
	31: {}, 32: {}, 33: {}, 34: {}, 35: {}, 37: {}, 38: {},
	41: {}, 42: {}, 43: {}, 44: {}, 45: {}, 46: {}, 47: {}, 48: {}, 49: {},
	51: {}, 53: {}, 54: {}, 55: {},
	61: {}, 62: {}, 63: {}, 64: {}, 65: {}, 66: {}, 67: {}, 68: {}, 69: {},
	71: {}, 73: {}, 74: {}, 75: {}, 77: {}, 79: {},
	81: {}, 82: {}, 83: {}, 84: {}, 85: {}, 86: {}, 87: {}, 88: {}, 89: {},
	91: {}, 92: {}, 93: {}, 94: {}, 95: {}, 96: {}, 97: {}, 98: {}, 99: {},
}

// ValidatePhone validates and formats a phone number, including international numbers.
// Returns: ok, formatted, errorMessage
func ValidatePhone(numero string) (string, error) {
	// keep only digits and '+'
	re := regexp.MustCompile(`[^\d\+]`)
	cleaned := re.ReplaceAllString(numero, "")

	if strings.HasPrefix(cleaned, "+55") || strings.HasPrefix(cleaned, "55") || strings.HasPrefix(cleaned, "0") {
		return validateBrazilian(cleaned)
	}
	if strings.HasPrefix(cleaned, "+") {
		// drop '+' for international validator
		return validateInternational(cleaned[1:])
	}
	return validateInternational(cleaned)
}

func validateBrazilian(numero string) (string, error) {
	switch {
	case strings.HasPrefix(numero, "55"):
		numero = numero[2:]
	case strings.HasPrefix(numero, "+55"):
		numero = numero[3:]
	case strings.HasPrefix(numero, "0"):
		numero = numero[1:]
	}

	if len(numero) != 10 && len(numero) != 11 {
		return "", errors.New("número inválido. Deve conter 10 ou 11 dígitos incluindo o DDD")
	}

	if len(numero) < 2 {
		return "", errors.New("número inválido. DDD ausente")
	}
	ddd, err := strconv.Atoi(numero[:2])
	if err != nil {
		return "", errors.New("DDD inválido")
	}
	if _, ok := brValidDDDs[ddd]; !ok {
		return "", errors.New("DDD inválido. Não está na lista de DDDs válidos")
	}

	if len(numero) == 10 {
		// if the third digit is 8 or 9, add a leading '9' after DDD
		if numero[2] == '8' || numero[2] == '9' {
			numero = numero[:2] + "9" + numero[2:]
		}
	}

	return "55" + numero, nil
}

func validateInternational(numero string) (string, error) {
	if len(numero) < 11 {
		return "", errors.New("número internacional inválido. Deve conter pelo menos 11 dígitos")
	}

	// choose 3-digit country code if first 3 are digits; otherwise 2
	codigoPais := ""
	if len(numero) >= 3 && isDigits(numero[:3]) {
		codigoPais = numero[:3]
	} else if len(numero) >= 2 {
		codigoPais = numero[:2]
	} else {
		return "", errors.New("código do país inválido")
	}

	if !isDigits(codigoPais) || len(codigoPais) > 3 {
		return "", errors.New("código do país inválido")
	}
	if strings.HasPrefix(codigoPais, "0") {
		return "", errors.New("código do pais não pode começar com zero")
	}

	return "+" + numero, nil
}

func isDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
