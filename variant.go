package bandit

import (
	"errors"
)

//An individual variant of an experiement.
type Variant interface {
	RewardSum() float64
	RewardSquareSum() float64
	ObservationCount() int
	Observe(reward float64) (Variant, error)
}

var OutOfRangeReward = errors.New("Reward is out of range [0, 1)")

// variant implements the Variant interface
var _ Variant = variant{}

// In memory implementation of a variant.
type variant struct {
	rewardSum        float64
	rewardSquareSum  float64
	observationCount int "oc"
}

func (v variant) RewardSum() float64 {
	return v.rewardSum
}

func (v variant) ObservationCount() int {
	return v.observationCount
}

func (v variant) RewardSquareSum() float64 {
	return v.rewardSquareSum
}
func NewVariant() Variant {
	return &variant{}
}

func (v variant) Observe(reward float64) (Variant, error) {
	if err := checkReward(reward); err != nil {
		return v, err
	}
	// Updating the variant estimates
	v.rewardSum += reward
	v.rewardSquareSum += reward * reward
	v.observationCount += 1
	return v, nil
}

func checkReward(reward float64) error {
	if reward < 0 {
		return OutOfRangeReward
	}

	if reward > 1 {
		return OutOfRangeReward
	}

	return nil
}
