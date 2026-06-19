package color

import "fmt"

const (
	ansiReset = "\033[0m"

	ansiBlack           = "\033[0;30m"
	ansiRed             = "\033[0;31m"
	ansiGreen           = "\033[0;32m"
	ansiYellow          = "\033[0;33m"
	ansiBlue            = "\033[0;34m"
	ansiMagenta         = "\033[0;35m"
	ansiCyan            = "\033[0;36m"
	ansiWhite           = "\033[0;37m"
	ansiGray            = "\033[90m"
	ansiPastelYellow    = "\033[38;2;255;245;150m"
	ansi256PastelYellow = "\033[38;5;229m"

	ansiBoldBlack   = "\033[1;30m"
	ansiBoldRed     = "\033[1;31m"
	ansiBoldGreen   = "\033[1;32m"
	ansiBoldYellow  = "\033[1;33m"
	ansiBoldBlue    = "\033[1;34m"
	ansiBoldMagenta = "\033[1;35m"
	ansiBoldCyan    = "\033[1;36m"
	ansiBoldWhite   = "\033[1;37m"

	ansiBold      = "\033[1m"
	ansiUnderline = "\033[4m"

	ansiBgRed   = "\033[41m"
	ansiBgGreen = "\033[42m"
)

func color(code string, a any) string {
	return code + fmt.Sprint(a) + ansiReset
}

func Black(a any) string   { return color(ansiBlack, a) }
func Red(a any) string     { return color(ansiRed, a) }
func Green(a any) string   { return color(ansiGreen, a) }
func Yellow(a any) string  { return color(ansiYellow, a) }
func Blue(a any) string    { return color(ansiBlue, a) }
func Magenta(a any) string { return color(ansiMagenta, a) }
func Cyan(a any) string    { return color(ansiCyan, a) }
func White(a any) string   { return color(ansiWhite, a) }
func Gray(a any) string    { return color(ansiGray, a) }

func BoldBlack(a any) string   { return color(ansiBoldBlack, a) }
func BoldRed(a any) string     { return color(ansiBoldRed, a) }
func BoldGreen(a any) string   { return color(ansiBoldGreen, a) }
func BoldYellow(a any) string  { return color(ansiBoldYellow, a) }
func BoldBlue(a any) string    { return color(ansiBoldBlue, a) }
func BoldMagenta(a any) string { return color(ansiBoldMagenta, a) }
func BoldCyan(a any) string    { return color(ansiBoldCyan, a) }
func BoldWhite(a any) string   { return color(ansiBoldWhite, a) }

func Bold(a any) string      { return color(ansiBold, a) }
func Underline(a any) string { return color(ansiUnderline, a) }

func ErrorLabel(a any) string   { return color(ansiBoldRed, a) }
func WarningLabel(a any) string { return color(ansiBoldYellow, a) }
func NoteLabel(a any) string    { return color(ansiBoldBlue, a) }

func Highlight(a any) string {
	return ansiUnderline + ansiBold + fmt.Sprint(a) + ansiReset
}

func PastelYellow(a any) string    { return color(ansiPastelYellow, a) }
func PastelYellow256(a any) string { return color(ansi256PastelYellow, a) }
