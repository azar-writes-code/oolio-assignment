package coupon

type Coupon struct {
	Code      string
	FilesMask byte
}

type CouponRepository interface {
	GetByCode(code string) (*Coupon, error)
}

type CouponService interface {
	Validate(code string) (bool, error)
}
