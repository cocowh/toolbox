package logger

import "fmt"

const (
	RedColor    = "\033[31m"
	GreenColor  = "\033[32m"
	YellowColor = "\033[33m"
	BlueColor   = "\033[34m"
	CyanColor   = "\033[36m"
	ResetColor  = "\033[0m"
)

func Red(msg string) string {
	return fmt.Sprintf("%s%s%s", RedColor, msg, ResetColor)
}

func Green(msg string) string {
	return fmt.Sprintf("%s%s%s", GreenColor, msg, ResetColor)
}

func Yellow(msg string) string {
	return fmt.Sprintf("%s%s%s", YellowColor, msg, ResetColor)
}

func Blue(msg string) string {
	return fmt.Sprintf("%s%s%s", BlueColor, msg, ResetColor)
}

func Cyan(msg string) string {
	return fmt.Sprintf("%s%s%s", CyanColor, msg, ResetColor)
}
