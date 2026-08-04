package money

import "testing"

func TestParseRupees(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    int64
		wantErr bool
	}{
		{"integer rupees", "250", 25000, false},
		{"with paise", "250.50", 25050, false},
		{"single decimal", "250.5", 25050, false},
		{"single paise padded", "250.05", 25005, false},
		{"zero", "0", 0, false},
		{"zero decimal", "0.00", 0, false},
		{"leading decimal", ".50", 50, false},
		{"whitespace trimmed", "  100.10  ", 10010, false},
		{"empty", "", 0, true},
		{"negative", "-5", 0, true},
		{"too many decimals", "1.2.3", 0, true},
		{"non numeric", "abc", 0, true},
		{"three decimals", "1.234", 0, true},
		{"trailing dot", "1.", 0, true},
		{"starts with dot empty", ".", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRupees(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseRupees(%q) expected error, got %d", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseRupees(%q) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseRupees(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestFormatPaise(t *testing.T) {
	tests := []struct {
		in   int64
		want string
	}{
		{125050, "1250.50"},
		{0, "0.00"},
		{5, "0.05"},
		{50, "0.50"},
		{100000, "1000.00"},
		{-125050, "-1250.50"},
	}
	for _, tt := range tests {
		if got := FormatPaise(tt.in); got != tt.want {
			t.Errorf("FormatPaise(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	inputs := []string{"0.00", "0.01", "0.50", "1.00", "250.50", "99999999.99"}
	for _, in := range inputs {
		paise, err := ParseRupees(in)
		if err != nil {
			t.Fatalf("ParseRupees(%q): %v", in, err)
		}
		if out := FormatPaise(paise); out != in {
			t.Errorf("round trip %q -> %q", in, out)
		}
	}
}
