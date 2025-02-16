package utils

import (
	"testing"

	"github.com/skykosiner/gobutar/pkg/items"
)

func TestFormatFloat(t *testing.T) {
	inputs := []struct {
		input    float64
		expected string
	}{
		{3.20, "3.20"},
		{2.92, "2.92"},
		{253.24, "253.24"},
	}

	for _, input := range inputs {
		if res := FormatFloat(input.input); res != input.expected {
			t.Logf("Expected %s for %f, got %s.", input.expected, input.input, res)
			t.Fail()
		}
	}
}

func TestFormatRecurring(t *testing.T) {
	inputs := []struct {
		input    items.Recurring
		expected string
	}{
		{items.No, "One Time"},
		{items.Daily, "Daily"},
		{items.Weekly, "Weekly"},
		{items.Monthly, "Monthly"},
		{items.Yearly, "Yearly"},
	}

	for _, input := range inputs {
		if res := FormatRecurring(input.input); res != input.expected {
			t.Logf("Expected %s for %s, got %s.", input.expected, input.input, res)
			t.Fail()
		}
	}
}
