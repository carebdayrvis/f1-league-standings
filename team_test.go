package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReversedScores(t *testing.T) {

	t.Run("5 results", func(t *testing.T) {
		e := []int{
			5,
			4,
			3,
			2,
			1,
		}

		assert.Equal(t, e, ReversedScores(5))
	})

	t.Run("20 results", func(t *testing.T) {
		e := []int{
			20,
			19,
			18,
			17,
			16,
			15,
			14,
			13,
			12,
			11,
			10,
			9,
			8,
			7,
			6,
			5,
			4,
			3,
			2,
			1,
		}

		assert.Equal(t, e, ReversedScores(20))
	})
}
