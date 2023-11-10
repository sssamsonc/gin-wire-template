package common_util

func GetStringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func GetIntValue(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

func GetBoolValue(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}
