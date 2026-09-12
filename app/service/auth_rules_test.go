package service

import (
	"testing"

	"api-students/app/model"
)

func TestValidateRegister(t *testing.T) {
	cases := []struct {
		name    string
		req     model.RegisterRequest
		wantErr []string
	}{
		{
			name:    "semua kosong",
			req:     model.RegisterRequest{Username: "", Email: "", Password: ""},
			wantErr: []string{"username", "email", "password"},
		},
		{
			name:    "valid",
			req:     model.RegisterRequest{Username: "budi123", Email: "budi@mail.com", Password: "rahasia123"},
			wantErr: []string{},
		},
		{
			name:    "username terlalu pendek",
			req:     model.RegisterRequest{Username: "ab", Email: "budi@mail.com", Password: "rahasia123"},
			wantErr: []string{"username"},
		},
		{
			name:    "email tidak valid",
			req:     model.RegisterRequest{Username: "budi123", Email: "bukan-email", Password: "rahasia123"},
			wantErr: []string{"email"},
		},
		{
			name:    "password terlalu pendek",
			req:     model.RegisterRequest{Username: "budi123", Email: "budi@mail.com", Password: "abc123"},
			wantErr: []string{"password"},
		},
		{
			name:    "password tanpa angka",
			req:     model.RegisterRequest{Username: "budi123", Email: "budi@mail.com", Password: "rahasiasaja"},
			wantErr: []string{"password"},
		},
		{
			name:    "password terlalu umum",
			req:     model.RegisterRequest{Username: "budi123", Email: "budi@mail.com", Password: "password1"},
			wantErr: []string{"password"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateRegister(tc.req)
			if len(errs) != len(tc.wantErr) {
				t.Errorf("dapat %d error, harap %d (%v)", len(errs), len(tc.wantErr), errs)
			}
			for _, field := range tc.wantErr {
				if _, ok := errs[field]; !ok {
					t.Errorf("field %q seharusnya error, tapi tidak", field)
				}
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	cases := []struct {
		name    string
		req     model.LoginRequest
		wantErr int
	}{
		{"semua kosong", model.LoginRequest{Username: "", Password: ""}, 2},
		{"valid", model.LoginRequest{Username: "budi123", Password: "apapun"}, 0},
		{"password kosong", model.LoginRequest{Username: "budi123", Password: ""}, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateLogin(tc.req)
			if len(errs) != tc.wantErr {
				t.Errorf("dapat %d error, harap %d (%v)", len(errs), tc.wantErr, errs)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	cases := []struct {
		username string
		want     bool
	}{
		{"budi123", true},
		{"budi.santoso", true},
		{"budi_santoso", true},
		{"budi santoso", false}, // ada spasi
		{"budi@santoso", false}, // ada karakter aneh
	}

	for _, tc := range cases {
		if got := isValidUsername(tc.username); got != tc.want {
			t.Errorf("isValidUsername(%q) = %v, harap %v", tc.username, got, tc.want)
		}
	}
}