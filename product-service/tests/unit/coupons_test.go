package unit

import (
	"testing"

	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/coupon"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/services/coupons"
)

type mockCouponRepo struct {
	data map[string]uint8
}

func (m *mockCouponRepo) GetByCode(code string) (*coupon.Coupon, error) {
	mask, ok := m.data[code]
	if !ok {
		return nil, nil
	}
	return &coupon.Coupon{Code: code, FilesMask: mask}, nil
}

func TestCouponService_Validate(t *testing.T) {
	repo := &mockCouponRepo{
		data: map[string]uint8{
			"VALID123":    0x03, // 011 -> in file 1 and 2
			"VALID456":    0x07, // 111 -> in all 3 files
			"INVALID1":    0x01, // 001 -> only in file 1
			"INVALID2":    0x04, // 100 -> only in file 3
			"SHORT":       0x03, // Too short
			"VERYLONG111": 0x03, // Too long (11 chars)
		},
	}
	service := coupons.NewService(repo)

	tests := []struct {
		code    string
		want    bool
		wantErr bool
	}{
		{"VALID123", true, false},
		{"VALID456", true, false},
		{"INVALID1", false, false},
		{"INVALID2", false, false},
		{"SHORT", false, false},
		{"VERYLONG111", false, false},
		{"NONEXISTENT", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := service.Validate(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Validate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
