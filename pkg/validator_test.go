package validator

import (
	"errors"
	"testing"
)

func TestValidateTodo(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		wantErr     bool
	}{
		{
			name:        "valid todo",
			title:       "Complete math homework",
			description: "Solve tasks 5 and 6 from mathematics book.",
			wantErr:     false,
		},
		{
			name:        "empty title",
			title:       "",
			description: "To make healthy breakfast to my little brother",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTodo(tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTodo() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrInvalidTitle) {
				t.Errorf("ValidateTodo() wrong error type: %v, want: %v", err, ErrInvalidTitle)
			}
		})
	}
}
