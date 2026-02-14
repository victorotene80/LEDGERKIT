package money

import (
	"math"
	"strings"

	lkerr "github.com/Helen-projects/LEDGERKIT/ledgerkit/errors"
)


type Money struct {
	asset AssetCode
	scale uint8
	amountMinor int64
}

func New(asset AssetCode, scale uint8, amountMinor int64) (Money, error) {
	asset = AssetCode(strings.ToUpper(strings.TrimSpace(asset.String())))
	if !asset.Valid() {
		return Money{}, lkerr.ErrInvalidAssetCode
	}
	if scale > 18 {
		return Money{}, lkerr.ErrInvalidScale
	}
	return Money{asset: asset, scale: scale, amountMinor: amountMinor}, nil
}

func MustNew(asset AssetCode, scale uint8, amountMinor int64) Money {
	m, err := New(asset, scale, amountMinor)
	if err != nil {
		panic(err)
	}
	return m
}

func (m Money) Asset() AssetCode   { return m.asset }
func (m Money) Scale() uint8       { return m.scale }
func (m Money) Minor() int64       { return m.amountMinor }
func (m Money) IsZero() bool       { return m.amountMinor == 0 }
func (m Money) IsPositive() bool   { return m.amountMinor > 0 }
func (m Money) IsNegative() bool   { return m.amountMinor < 0 }
func (m Money) Abs() Money {
	if m.amountMinor >= 0 {
		return m
	}
	if m.amountMinor == math.MinInt64 {
		panic(lkerr.ErrOverflow)
	}
	m.amountMinor = -m.amountMinor
	return m
}

func (m Money) EnsureSameAsset(other Money) error {
	if m.asset != other.asset {
		return lkerr.ErrAssetMismatch
	}
	if m.scale != other.scale {
		return lkerr.ErrScaleMismatch
	}
	return nil
}

func (m Money) Add(other Money) (Money, error) {
	if err := m.EnsureSameAsset(other); err != nil {
		return Money{}, err
	}
	sum, ok := safeAddInt64(m.amountMinor, other.amountMinor)
	if !ok {
		return Money{}, lkerr.ErrOverflow
	}
	m.amountMinor = sum
	return m, nil
}

func (m Money) Sub(other Money) (Money, error) {
	if err := m.EnsureSameAsset(other); err != nil {
		return Money{}, err
	}
	diff, ok := safeSubInt64(m.amountMinor, other.amountMinor)
	if !ok {
		return Money{}, lkerr.ErrOverflow
	}
	m.amountMinor = diff
	return m, nil
}

func safeAddInt64(a, b int64) (int64, bool) {
	if (b > 0 && a > math.MaxInt64-b) || (b < 0 && a < math.MinInt64-b) {
		return 0, false
	}
	return a + b, true
}

func safeSubInt64(a, b int64) (int64, bool) {
	if b == math.MinInt64 {
		return 0, false
	}
	return safeAddInt64(a, -b)
}
