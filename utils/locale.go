package utils

func GetLanguage(locale string) string {
	switch locale {
	case "tr-TR":
		return "tr"
	case "en-GB":
		return "en"
	default:
		return "tr"
	}
}
