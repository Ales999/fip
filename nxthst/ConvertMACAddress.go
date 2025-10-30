package nxthst

import (
	"fmt"
	"regexp"
	"strings"
)

// convertMACAddress проверяет и при необходимости конвертирует MAC-адрес
func ConvertMACAddress(input string) string {

	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return input
	}

	// Регулярное выражение для проверки формата xx:xx:xx:xx:xx:xx
	macRegexHex := regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)
	if macRegexHex.MatchString(input) {
		// Удаляем все двоеточия и приводим к нижнему регистру
		cleaned := strings.ReplaceAll(input, ":", "")
		cleaned = strings.ToLower(cleaned)

		return macReformat(cleaned, input)
	}

	// Регулярное выражение для проверки формата 01xx.xxxx.xxxx.xx
	macRegexDhcp := regexp.MustCompile(`^01([0-9a-fA-F]{2}\.)([0-9a-fA-F]{4}\.){2}([0-9a-fA-F]{2})$`)
	if macRegexDhcp.MatchString(input) {
		cleaned := strings.ReplaceAll(input, ".", "")
		cleaned = strings.ToLower(cleaned)
		cleaned = cleaned[2:]

		return macReformat(cleaned, input)
	}
	// Если ничего не подошло, то вернем все как есть.
	return input

}

// macReformat - Делим входную строку на части по 4 символа и объединяем точками
func macReformat(cleaned string, original string) string {
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
	if !strings.EqualFold(original, ret) {
		fmt.Println("Convert", original, "-->", ret)
	}

	return ret
}
