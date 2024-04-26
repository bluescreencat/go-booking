package util_test

import (
	"booking/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("the password after hash password again should not get the same as previous value", func(t *testing.T) {
		util := util.New()
		password := "thisispassword"
		round := 10
		got1, _ := util.HashPassword(password, round)
		got2, _ := util.HashPassword(password, round)
		assert.NotEqualValues(t, got1, got2, "password should not get the same value")
	})
}

func TestComparePassword(t *testing.T) {
	t.Run("the password should comparable", func(t *testing.T) {
		util := util.New()
		password := "thisispassword"
		round := 10
		hashPassword, _ := util.HashPassword(password, round)
		if err := util.ComparePassword(hashPassword, password); err != nil {
			t.Error("the password should comparable")
		}
	})
}
