package bandit

import (
	"testing"
)

func TestRoundIndex(t *testing.T) {
	var tests = []struct {
		in  []Variant
		out int
	}{
		{[]Variant{}, 0},
		{[]Variant{variant{}}, 0},
		{[]Variant{variant{0, 0, 1}}, 1},
		{[]Variant{variant{0, 0, 1}, variant{0, 0, 2}}, 3},
	}
	for _, v := range tests {
		if out := RoundIndex(v.in); out != v.out {
			t.Errorf("Unexpecte round index for %v. Expecting %d but found %d", v, v.out, out)
		}
	}
}

func TestVariantMean(t *testing.T) {
	var tests = []struct {
		in  variant
		out float64
	}{
		{variant{}, 0},
		{variant{0, 0, 1}, 0},
		{variant{1, 1, 1}, 1},
		{variant{1, 1, 2}, 0.5},
	}

	for _, v := range tests {
		if mean := Mean(v.in); mean != v.out {
			t.Errorf("Unexpected variant mean for %v. Expecting %f but found %f", v, v.out, mean)
		}
	}
}

func TestObservedCount(t *testing.T) {
	var tests = []struct {
		in  []Variant
		out int
	}{
		{[]Variant{}, 0},
		{[]Variant{variant{0, 0, 0}}, 0},
		{[]Variant{variant{0, 0, 0}, variant{0, 0, 0}}, 0},
		{[]Variant{variant{0, 0, 1}, variant{0, 0, 0}}, 1},
		{[]Variant{variant{0, 0, 0}, variant{0, 0, 2}}, 1},
		{[]Variant{variant{0, 0, 3}, variant{0, 0, 2}}, 2},
		{[]Variant{variant{0, 0, 3}, variant{0, 0, 0}, variant{0, 0, 2}}, 2},
	}

	for _, v := range tests {
		if count := ObservedCount(v.in); count != v.out {
			t.Errorf("Unexpected observed count for %v. expecting %d but found %d.", v, v.out, count)
		}
	}
}

func TestTwiceObservedCount(t *testing.T) {
	var tests = []struct {
		in  []Variant
		out int
	}{
		{[]Variant{}, 0},
		{[]Variant{variant{0, 0, 0}}, 0},
		{[]Variant{variant{0, 0, 0}, variant{0, 0, 0}}, 0},
		{[]Variant{variant{0, 0, 1}, variant{0, 0, 0}}, 0},
		{[]Variant{variant{0, 0, 0}, variant{0, 0, 2}}, 1},
		{[]Variant{variant{0, 0, 1}, variant{0, 0, 2}}, 1},
		{[]Variant{variant{0, 0, 3}, variant{0, 0, 2}}, 2},
		{[]Variant{variant{0, 0, 3}, variant{0, 0, 0}, variant{0, 0, 2}}, 2},
	}

	for _, v := range tests {
		if count := TwiceObservedCount(v.in); count != v.out {
			t.Errorf("Unexpected twice observed count for %v. expecting %d but found %d.", v, v.out, count)
		}
	}
}

func TestEpsilonGreedy(t *testing.T) {
	// In these tests we make use of the following facts about the function itself
	// if RoundIndex(variants) == 0 then we return a random variant
	// if epsilon < 0 then we return a random variant
	// if epsilon >= 1 AND RoundIndex(variants) == 0 then we force the greedest
	// If these assumptions are changed then these tests might fail.

	var tests = []struct {
		epsilon float64
		in      []Variant
		out     Variant
	}{
		{0, []Variant{variant{}}, variant{}},                                  //random
		{1, []Variant{variant{0, 0, 1}}, variant{0, 0, 1}},                    //random
		{-1, []Variant{variant{0, 0, 1}}, variant{0, 0, 1}},                   // When we have a single variant always return it.
		{-1, []Variant{variant{0, 0, 0}, variant{1, 1, 1}}, variant{1, 1, 1}}, // the one with the greater mean
		{-1, []Variant{variant{1, 1, 1}, variant{}}, variant{1, 1, 1}},        // the one with the greater mean
	}

	// Check that two variants are equal within a certain delta to account for floating point arithmetic.
	checkEquals := func(a, b Variant) bool {
		delta := 0.000001
		return (a.RewardSum()-b.RewardSum() < delta) &&
			(a.RewardSquareSum()-b.RewardSquareSum() < delta) &&
			a.ObservationCount() == b.ObservationCount()
	}

	for _, v := range tests {
		if variant := EpsilonGreedy(v.epsilon, v.in); !checkEquals(v.out, variant) {
			t.Errorf("Unexpected variant selected for %v. Expecting %v found %v", v.in, v.out, variant)
		}
	}

	//Test the case where we provide no []Variants to select from
	testValues := []float64{0, 1, -1, 0.5}
	for _, epsilon := range testValues {
		if v := EpsilonGreedy(epsilon, []Variant{}); v != nil {
			t.Errorf("epsilon = %d => Expecting the method to return nil when no variants are provided but found %v", v)
		}
	}
}
