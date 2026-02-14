package money

import (
	"testing"

	lkerr "github.com/victorotene80/LEDGERKIT/ledgerkit/errors"
)

func TestNewMoney_Valid(t *testing.T) {
	m, err := New(AssetCode("NGN"), 2, 1050)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if m.Asset() != "NGN" {
		t.Fatalf("expected NGN, got %s", m.Asset())
	}
	if m.Scale() != 2 {
		t.Fatalf("expected scale 2, got %d", m.Scale())
	}
	if m.Minor() != 1050 {
		t.Fatalf("expected 1050, got %d", m.Minor())
	}
}

func TestNewMoney_InvalidAsset(t *testing.T) {
	_, err := New(AssetCode("ngn!"), 2, 100)
	if err != lkerr.ErrInvalidAssetCode {
		t.Fatalf("expected ErrInvalidAssetCode, got %v", err)
	}
}

func TestNewMoney_InvalidScale(t *testing.T) {
	_, err := New(AssetCode("USD"), 50, 100)
	if err != lkerr.ErrInvalidScale {
		t.Fatalf("expected ErrInvalidScale, got %v", err)
	}
}

func TestEnsureSameAsset(t *testing.T) {
	a := MustNew("NGN", 2, 100)
	b := MustNew("USD", 2, 100)

	if err := a.EnsureSameAsset(b); err != lkerr.ErrAssetMismatch {
		t.Fatalf("expected ErrAssetMismatch, got %v", err)
	}
}

func TestEnsureSameScale(t *testing.T) {
	a := MustNew("NGN", 2, 100)
	b := MustNew("NGN", 0, 100)

	if err := a.EnsureSameAsset(b); err != lkerr.ErrScaleMismatch {
		t.Fatalf("expected ErrScaleMismatch, got %v", err)
	}
}

func TestAdd_Sub(t *testing.T) {
	a := MustNew("NGN", 2, 1000)
	b := MustNew("NGN", 2, 250)

	sum, err := a.Add(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.Minor() != 1250 {
		t.Fatalf("expected 1250, got %d", sum.Minor())
	}

	diff, err := sum.Sub(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff.Minor() != 1000 {
		t.Fatalf("expected 1000, got %d", diff.Minor())
	}
}
