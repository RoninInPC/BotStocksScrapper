package markdown

import "fmt"

// функции для Markdown
func ToBold(str string) string {
	return fmt.Sprintf("*%s*", str)
}

func ToItalic(str string) string {
	return fmt.Sprintf("_%s_", str)
}
