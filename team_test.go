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

func TestFormatModifier(t *testing.T) {

	t.Run("both", func(t *testing.T) {
		r := Race{
			FastestLapPoint:  true,
			LastInQualyPoint: true,
		}

		assert.Equal(t, "(+1 for fastest lap, and last in qualifying)", FormatModifier(r))

	})

	t.Run("neither", func(t *testing.T) {
		r := Race{
			FastestLapPoint:  false,
			LastInQualyPoint: false,
		}

		assert.Equal(t, "", FormatModifier(r))
	})

	t.Run("fastest", func(t *testing.T) {
		r := Race{
			FastestLapPoint:  true,
			LastInQualyPoint: false,
		}

		assert.Equal(t, "(+1 for fastest lap)", FormatModifier(r))
	})

	t.Run("qualy", func(t *testing.T) {
		r := Race{
			FastestLapPoint:  false,
			LastInQualyPoint: true,
		}

		assert.Equal(t, "(+1 for last in qualifying)", FormatModifier(r))
	})
}
