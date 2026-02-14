package ledgerkit

import "testing"

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Fatal("Version should not be empty")
	}
}
