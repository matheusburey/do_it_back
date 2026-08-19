package pkg

import "testing"

func TestNotBlank(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"empty", "", false},
		{"only spaces", "   ", false},
		{"has content", "hello", true},
		{"padded content", "  hello  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotBlank(tt.value); got != tt.want {
				t.Errorf("NotBlank(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		length int
		want   bool
	}{
		{"shorter", "abc", 5, false},
		{"exact", "abcde", 5, true},
		{"longer", "abcdef", 5, true},
		{"multibyte runes count as one", "áéíóú", 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinLength(tt.value, tt.length); got != tt.want {
				t.Errorf("MinLength(%q, %d) = %v, want %v", tt.value, tt.length, got, tt.want)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		length int
		want   bool
	}{
		{"shorter", "abc", 5, true},
		{"exact", "abcde", 5, true},
		{"longer", "abcdef", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxLength(tt.value, tt.length); got != tt.want {
				t.Errorf("MaxLength(%q, %d) = %v, want %v", tt.value, tt.length, got, tt.want)
			}
		})
	}
}

func TestIsEmail(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid email", "user@example.com", true},
		{"valid with subdomain", "user@mail.example.com", true},
		{"missing at", "userexample.com", false},
		{"missing domain", "user@", false},
		{"missing local part", "@example.com", false},
		{"double at", "user@@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmail(tt.value); got != tt.want {
				t.Errorf("IsEmail(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsPassword(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"meets all rules", "Abcdef1!", true},
		{"missing upper", "abcdef1!", false},
		{"missing lower", "ABCDEF1!", false},
		{"missing number", "Abcdefg!", false},
		{"missing symbol", "Abcdefg1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPassword(tt.value); got != tt.want {
				t.Errorf("IsPassword(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestEvaluator_CheckField(t *testing.T) {
	var eval Evaluator

	eval.CheckField(true, "name", "should not appear")
	eval.CheckField(false, "name", "name is required")

	if len(eval) != 1 {
		t.Fatalf("expected 1 field error, got %d", len(eval))
	}

	if msg := eval["name"]; msg != "name is required" {
		t.Errorf("expected message %q, got %q", "name is required", msg)
	}
}

func TestEvaluator_KeepsFirstErrorPerField(t *testing.T) {
	var eval Evaluator

	eval.AddFieldError("email", "email is invalid")
	eval.AddFieldError("email", "email is too short")

	if msg := eval["email"]; msg != "email is invalid" {
		t.Errorf("expected first error to be kept, got %q", msg)
	}
}
