package bandit

import (
	"testing"
)

func TestObserve(t *testing.T) {
	var variantsTest = []struct {
		reward float64
		in     variant
		out    variant
	}{
		{0, variant{}, variant{0, 0, 1}},
		{0.1, variant{0, 0, 1}, variant{0.1, 0.1 * 0.1, 2}},
		{0.9, variant{0.1, 0.5, 2}, variant{1, 1.31, 3}},
	}

	// Check that two variants are equal within a certain delta to account for floating point arithmetic.
	checkEquals := func(a, b Variant) bool {
		delta := 0.000001
		return (a.RewardSum()-b.RewardSum() < delta) &&
			(a.RewardSquareSum()-b.RewardSquareSum() < delta) &&
			a.ObservationCount() == b.ObservationCount()
	}

	for _, v := range variantsTest {
		out, err := v.in.Observe(v.reward)
		if err != nil {
			t.Errorf("Unexpected error %v, ", err)
		}
		if !checkEquals(out, v.out) {
			t.Errorf("Unexpected result for %v. out=%v", v, out)
		}
	}

}

func TestObserveWrongRange(t *testing.T) {
	v := variant{}
	_, err := v.Observe(-0.00001)
	if err != OutOfRangeReward {
		t.Errorf("Expecting to fail with an out of range reward error for %f", -0.00001)
	}

	_, err = v.Observe(1.0001)
	if err != OutOfRangeReward {
		t.Errorf("Expecting to fail with an out of range reward error for %f", 1.0001)
	}
}
