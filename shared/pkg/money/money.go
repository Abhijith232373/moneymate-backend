package money

import (
	"fmt"
	"strconv"
	"strings"
)

const maxRupees = (1 << 62) / 100 // plenty for any real wallet

// ParseRupees converts a user-supplied string like "250.50" or "250" into paise (int64).
// Returns an error for empty, negative, non-numeric, or more-than-2-decimal input.
func ParseRupees(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("amount is required")
	}
	if strings.HasPrefix(s, "-") {
		return 0, fmt.Errorf("amount must be a positive value")
	}

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid amount %q", s)
	}

	rupees := parts[0]
	if rupees == "" {
		rupees = "0"
	}
	if !allDigits(rupees) {
		return 0, fmt.Errorf("invalid amount %q", s)
	}

	paise := "00"
	if len(parts) == 2 {
		p := parts[1]
		if p == "" || !allDigits(p) || len(p) > 2 {
			return 0, fmt.Errorf("invalid amount %q", s)
		}
		paise = p
		for len(paise) < 2 {
			paise += "0"
		}
	}

	r, err := strconv.ParseInt(rupees, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount %q", s)
	}
	p, err := strconv.ParseInt(paise, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount %q", s)
	}
	if r > maxRupees {
		return 0, fmt.Errorf("amount too large")
	}
	return r*100 + p, nil
}

// FormatPaise converts paise to a rupees string, e.g. 125050 → "1250.50".
func FormatPaise(p int64) string {
	sign := ""
	if p < 0 {
		sign = "-"
		p = -p
	}
	return fmt.Sprintf("%s%d.%02d", sign, p/100, p%100)
}

func allDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
