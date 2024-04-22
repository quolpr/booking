package days

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDaysBetween(t *testing.T) {
	t.Parallel()

	t.Run("Works correctly", func(t *testing.T) {
		t.Parallel()
		res, err := DaysBetween(Date(2021, 1, 2), Date(2021, 1, 3))

		assert.Nil(t, err)

		assert.Equal(t, []time.Time{Date(2021, 1, 2), Date(2021, 1, 3)}, res)
	})

	t.Run("Works with only one day when from == to", func(t *testing.T) {
		t.Parallel()
		res, err := DaysBetween(Date(2021, 1, 2), Date(2021, 1, 2))

		assert.Nil(t, err)

		assert.Equal(t, []time.Time{Date(2021, 1, 2)}, res)
	})

	t.Run("Return error on wrong range", func(t *testing.T) {
		t.Parallel()

		res, err := DaysBetween(Date(2021, 1, 2), Date(2021, 1, 1))

		assert.Empty(t, res)

		assert.ErrorIs(t, err, ErrWrongRange)
	})
}
