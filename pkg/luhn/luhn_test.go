package luhn

import "testing"

func Test_Check(t *testing.T) {
    tests := []struct {
        input string
        want  bool
    }{
        {"5062821234567892", true},
        {"5062821734567892", false},
        {"123456789012345678901234567891", true},
        {"123456789012345678901234567890", false},
    }

    for _, tt := range tests {
        got := Check(tt.input)
        if got != tt.want {
            t.Errorf("Check(%q) = %v, want %v", tt.input, got, tt.want)
        }
    }
}

func Test_revers(t *testing.T) {
    tests := []struct {
        input string
        want  string
    }{
        {"1234567890", "0987654321"},
        {"gopher", "rehpog"},
        {"", ""},
    }

    for _, tt := range tests {
        got := revers(tt.input)
        if got != tt.want {
            t.Errorf("revers(%q) = %q, want %q", tt.input, got, tt.want)
        }
    }
}
