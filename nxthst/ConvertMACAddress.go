package nxthst

import (
	"fmt"
	"regexp"
	"strings"
)

// convertMACAddress проверяет и при необходимости конвертирует MAC-адрес
func ConvertMACAddress(input string) string {
	// Регулярное выражение для проверки формата xx:xx:xx:xx:xx:xx
	macRegex := regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)
	if !macRegex.MatchString(input) {
		return input
	}

	// Удаляем все двоеточия и приводим к нижнему регистру
	cleaned := strings.ReplaceAll(input, ":", "")
	cleaned = strings.ToLower(cleaned)

	// Делим строку на части по 4 символа и объединяем точками
	var result strings.Builder
	for i := 0; i < len(cleaned); i += 4 {
		end := i + 4
		if end > len(cleaned) {
			end = len(cleaned)
		}
		result.WriteString(cleaned[i:end])
		if end != len(cleaned) {
			result.WriteRune('.')
		}
	}

	ret := result.String()
	// Если коррекция была произведена то выведем что изменили.
	if !strings.EqualFold(input, ret) {
		fmt.Println("Convert", input, "-->", ret)
	}

	return ret
}
