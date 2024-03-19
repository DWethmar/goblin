package game

import (
	"strings"
	"testing"
)

func TestValidateName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty name",
			args:    args{name: ""},
			wantErr: true,
		},
		{
			name:    "name too long",
			args:    args{name: strings.Repeat("a", 101)},
			wantErr: true,
		},
		{
			name:    "invalid name: number",
			args:    args{name: "123"},
			wantErr: true,
		},
		{
			name:    "invalid name: special character",
			args:    args{name: "abc!"},
			wantErr: true,
		},
		{
			name:    "invalid name: uppercase",
			args:    args{name: "ABC"},
			wantErr: true,
		},
		{
			name:    "valid name",
			args:    args{name: "abc"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
