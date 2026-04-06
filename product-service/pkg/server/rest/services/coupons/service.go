package coupons

import (
	"math/bits"

	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/coupon"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/utils/apperrors"
)

type service struct {
	repo coupon.CouponRepository
}

func NewService(repo coupon.CouponRepository) coupon.CouponService {
	return &service{repo: repo}
}

func (s *service) Validate(code string) (bool, error) {
	// Rule 1: Length between 8 and 10 characters
	if len(code) < 8 || len(code) > 10 {
		return false, nil
	}

	// Rule 2: Found in at least two files
	c, err := s.repo.GetByCode(code)
	if err != nil {
		return false, apperrors.NewInternal("failed to get coupon", err)
	}

	if c == nil {
		return false, apperrors.NewNotFound("coupon code %s not found", code)
	}

	// Count set bits in the mask
	count := bits.OnesCount8(c.FilesMask)

	return count >= 2, nil
}
