package helpers

func TranslateNewsType(newsType string) string {
	switch newsType {
	case "regular":
		return "Регулярная"
	case "breaking":
		return "Срочная"
	default:
		return "Неизвестно"
	}
}
