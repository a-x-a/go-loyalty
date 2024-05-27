package luhn

import (
	"regexp"
	"strconv"
)

// Check checks if a string is a valid order number using the Luhn algorithm.
//
// Parameters:
// - s: the string to be checked.
//
// Returns:
// - bool: true if the string is a valid order number, false otherwise.
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

// revers reverses a string.
//
// It takes a string as a parameter and returns the reversed string.
func revers(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; j > i; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}
