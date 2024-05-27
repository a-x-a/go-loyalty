package luhn

import (
	"regexp"
	"strconv"
)

func Check(s string) bool {
	digitsRegExp := regexp.MustCompile(`^\d+$`)
	if !digitsRegExp.MatchString(s) {
		return false
	}

	s = revers(s)

	sum := 0
	for i := 0; i < len(s); i++ {
		d, _ := strconv.Atoi(string(s[i]))
		if (i+1)%2 == 0 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}

		sum += d
	}

	return sum%10 == 0
}

func revers(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; j > i; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}
