package tests

import "testing"

func TestListAccounts(t *testing.T) {
	accts, err := c.api.ListAccounts(c.ctx)
	if err != nil {
		t.Errorf("failed listing accounts: %v", err)
		return
	}

	if len(accts) == 0 {
		t.Errorf("account list should at least be 1")
	}
}
