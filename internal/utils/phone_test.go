package utils

import (
	"strings"
	"testing"
)

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr string // substring match
	}{
		{
			name: "Brazil_with_plus55_11digits",
			in:   "+55 11 98877-6655",
			want: "5511988776655",
		},
		{
			name: "Brazil_without_plus_with_55_prefix",
			in:   "5511988776655",
			want: "5511988776655",
		},
		{
			name: "Brazil_with_trunk_zero",
			in:   "011988776655",
			want: "5511988776655",
		},
		{
			name: "Brazil_10_digits_add_9_when_third_is_8",
			in:   "+55 11 8812 3456",
			want: "5511988123456",
		},
		{
			name:    "Brazil_invalid_DDD",
			in:      "+55 10 912345678",
			wantErr: "DDD inválido",
		},
		{
			name:    "Brazil_invalid_length",
			in:      "+55 11 12345",
			wantErr: "10 ou 11 dígitos",
		},
		{
			name: "International_with_plus",
			in:   "+14155552671",
			want: "+14155552671",
		},
		{
			name: "International_without_plus",
			in:   "14155552671",
			want: "+14155552671",
		},
		{
			name:    "International_too_short",
			in:      "+1234567890",
			wantErr: "pelo menos 11",
		},
		{
			name:    "International_country_code_starts_with_zero",
			in:      "+01234567890",
			wantErr: "não pode começar com zero",
		},
		{
			name: "Brazil_cleans_formatting_characters",
			in:   "(+55) 11-98877-6655",
			want: "5511988776655",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := ValidatePhone(tc.in)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil (got=%q)", tc.wantErr, got)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("error %q does not contain %q", err.Error(), tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
