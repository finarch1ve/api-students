package service

import (
	"testing"

	"api-students/app/model"
)

func TestValidateCreate(t *testing.T) {
	cases := []struct {
		name    string
		req     model.CreateStudentRequest
		wantErr []string
	}{
		{
			name:    "semua kosong",
			req:     model.CreateStudentRequest{NIM: "", Name: "", Grade: 150},
			wantErr: []string{"nim", "name", "grade"},
		},
		{
			name:    "valid",
			req:     model.CreateStudentRequest{NIM: "123", Name: "Ayu", Grade: 85},
			wantErr: []string{},
		},
		{
			name:    "grade negatif",
			req:     model.CreateStudentRequest{NIM: "123", Name: "Ayu", Grade: -5},
			wantErr: []string{"grade"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateCreate(tc.req)
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

func TestValidateReplace(t *testing.T) {
	errs := ValidateReplace(model.ReplaceStudentRequest{NIM: "", Name: "", Grade: 50})
	if len(errs) != 2 {
		t.Errorf("harap 2 error (nim, name), dapat %d: %v", len(errs), errs)
	}
}

func TestApplyPatch(t *testing.T) {
	current := model.Student{ID: 1, NIM: "123", Name: "Ayu", Grade: 85, IsActive: true}
	inactive := false

	result, errs := ApplyPatch(current, model.PatchStudentRequest{IsActive: &inactive})

	if len(errs) != 0 {
		t.Fatalf("tidak seharusnya ada error: %v", errs)
	}
	if result.IsActive {
		t.Error("is_active seharusnya berubah menjadi false")
	}
	if result.Name != "Ayu" {
		t.Error("field yang tidak dikirim seharusnya tidak berubah")
	}
}

func TestCountTotalPages(t *testing.T) {
	cases := []struct{ total, limit, want int }{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
	}
	for _, tc := range cases {
		if got := CountTotalPages(tc.total, tc.limit); got != tc.want {
			t.Errorf("total=%d limit=%d: harap %d, dapat %d", tc.total, tc.limit, tc.want, got)
		}
	}
}