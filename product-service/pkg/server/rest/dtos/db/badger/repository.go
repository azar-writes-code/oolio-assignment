package badger

import (
	"errors"

	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/coupon"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/utils/apperrors"
	"github.com/dgraph-io/badger/v4"
)

type badgerRepository struct {
	db *badger.DB
}

// New creates a new Badger-backed CouponRepository.
func New(db *badger.DB) coupon.CouponRepository {
	return &badgerRepository{db: db}
}

func (r *badgerRepository) GetByCode(code string) (*coupon.Coupon, error) {
	var mask byte
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("promo:" + code))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			if len(val) > 0 {
				mask = val[0]
			}
			return nil
		})
	})

	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, apperrors.NewInternal("failed to get coupon from badger", err)
	}

	return &coupon.Coupon{
		Code:      code,
		FilesMask: mask,
	}, nil
}
