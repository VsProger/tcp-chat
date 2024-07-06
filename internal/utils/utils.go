package utils

func ProcessBackspaces(input string) string {
	var result []rune
	for _, r := range input {
		if r == '\b' {
			if len(result) > 0 {
				result = result[:len(result)-1]
			}
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
