package bandit

// Package bandit inclludes different strategies/algorithems for the
//stochastic multi-armed bandit problem 

import (
	"math"
	"math/rand"
)

// Given a group of variant the round index is the sum of the observation count.
func RoundIndex(variants []Variant) (roundIndex int) {
	for _, v := range variants {
		roundIndex += v.ObservationCount()
	}
	return
}

// Returns the number of variants that have been observed.
func ObservedCount(variants []Variant) (count int) {
	for _, v := range variants {
		if v.ObservationCount() > 0 {
			count++
		}
	}
	return
}

// Returns the number of variants that have been observed twice
func TwiceObservedCount(variants []Variant) (count int) {
	for _, v := range variants {
		if v.ObservationCount() > 1 {
			count++
		}
	}
	return
}

// Gets the reward mean associated to the specified variant
func Mean(v Variant) float64 {
	return v.RewardSum() / math.Max(float64(v.ObservationCount()), 1)
}

// Gets the reward standard deviation associated to the specified variant.
func Sigma(v Variant) float64 {
	if v.ObservationCount() == 0 {
		return math.NaN()
	}
	mean := Mean(v)
	variance := v.RewardSquareSum()/float64(v.ObservationCount()) - mean*mean

	return math.Sqrt(variance)
}

// The sum of the standard deviation for the given variants.
func SigmaSum(variants []Variant) (sum float64) {
	for _, v := range variants {
		if sigma := Sigma(v); !math.IsNaN(sigma) {
			sum += sigma
		}
	}
	return
}

// Greedy epsilon strategy.
// This is the most simple of strategies. Intuitively it basically always pulls the lever with
// the highest estimated mean, except when a randome lever is pulled with an epsilon frequency.
// The function will return nil if len(variants) == 0
func EpsilonGreedy(epsilon float64, variants []Variant) Variant {
	if len(variants) == 0 {
		return nil
	}
	if RoundIndex(variants) == 0 || rand.Float64() < epsilon {
		//Play a random variant
		return variants[rand.Intn(len(variants))]
	}

	return greatestMean(variants)
}

// Epsilon decreasing strategy is similar to the epsilon greed but the epsilon value
// decreases over time. This implementation provides a decreasing e0 / t  where e0 is a positive tuning
// parameter and t the current round index.
// The epsilon decreasing strategy is analysed in <i>Finite time analysis of the multiarmed
// bandit problem.</i> by Auer, Cesa-Bianchi and Fisher in <i>Machine Learning</i> (2002).
func EpsilonDecreasing(e0 float64, variants []Variant) Variant {
	epsilonZeroT := math.Min(1.0, e0/math.Max(float64(RoundIndex(variants)), 1))

	if rand.Float64() < epsilonZeroT {
		// randomom lever with epsilonZero frequency
		return variants[rand.Intn(len(variants))]
	}

	return greatestMean(variants)
}

// Heuristic policy for multi-armed bandit problem.
// Implementation is based on the algorithm presented by Auer et al in
// "Finite-Time Analysis of the Multiarmed Bandit Problem" (2002).
func Ucb1(variants []Variant) Variant {
	if len(variants) == 0 {
		return nil
	}
	maxUcb1 := math.Inf(-1)
	maxUcbIndex := -1
	//access the array in a randomized fashion so repeated invocations of this function don't return the same value.
	for _, index := range rand.Perm(len(variants)) {
		v := variants[index]
		if v.ObservationCount() == 0 {
			// If we have never played this variant play it now.
			return v
		}

		ucb1 := rank(v, RoundIndex(variants)) // the true mean.
		if ucb1 > maxUcb1 {
			maxUcb1 = ucb1
			maxUcbIndex = index
		}
	}
	return variants[maxUcbIndex]
}

// Rank is ucb1 value = mean + confidence bound
func rank(v Variant, totalObservations int) float64 {
	return Mean(v) + math.Sqrt(2.0*math.Log(float64(totalObservations))/float64(v.ObservationCount()))
}

//Play the variant with the greatest mean.
func greatestMean(variants []Variant) (result Variant) {
	maxMean := math.Inf(-1) //Set initial max mean to -Infinity
	for _, v := range variants {
		if mean := Mean(v); mean > maxMean {
			maxMean = mean
			result = v
		}
	}
	return
}
