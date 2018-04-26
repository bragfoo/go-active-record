package main

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func Test_formatName(t *testing.T) {
	type args struct {
		str string
		cfl bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"aaa", args{"aaa", true}, "Aaa"},
		{"Aaa", args{"Aaa", true}, "Aaa"},
		{"Aaa", args{"Aaa", false}, "aaa"},
		{"Aaa", args{"aaa", false}, "aaa"},
		{"a_a_a", args{"a_a_a", false}, "aAA"},
		{"a_a_", args{"a_a_", true}, "AA"},
		{"a_aa", args{"a_aa", true}, "AAa"},
		{"aa_a", args{"aa_a", true}, "AaA"},
		{"_a_a", args{"_a_a", true}, "AA"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatName(tt.args.str, tt.args.cfl); got != tt.want {
				t.Errorf("formatName() = %v, want %v", got, tt.want)
			}
		})
	}
}
