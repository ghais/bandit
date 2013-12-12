bandit
======
[![Build Status](https://travis-ci.org/ghais/bandit.png?branch=master)](https://travis-ci.org/ghais/bandit)
[![Coverage Status](https://coveralls.io/repos/ghais/bandit/badge.png)](https://coveralls.io/r/ghais/bandit)

For more information see [Multiarmed Bandit](http://en.wikipedia.org/wiki/Multi-armed_bandit)

Types
-----
```go
type Variant interface {
    RewardSum() float64 
	RewardSquareSum() float64
	ObservationCount() int
	Observe(reward float64) (Variant, error)
}
```
All the algorithm implementations operate on an array of variants. An in-memory Variant implementation 
is provided using bandit.NewVariant(). However you can implement your own persistent variant. 

Internally we have implemented persistent variants over Mongodb, Cassandra and Leveldb. We would probably open source those at a later date.

The in-memory variant implementation is immutable, and therefore when calling the Observe method
you should re-assign it:
```go
v := bandit.NewVaraint()
v, _ = v.Observe(0.5)
```

Algorithms and Examples
----------
* *UCB1*:
Fina the machine that maximizes the rank (mean + confidence bound).
Implementation is based on the algorithm presented by Auer et al in
"Finite-Time Analysis of the Multiarmed Bandit Problem" (2002).
```go
variants := []Variant{bandit.NewVariant(), bandit.NewVariant(), bandit.NewVariant()}
v, _ := bandit.Ucb1(variants)
```

* *EpsilonDecreasing*:
Epsilon decreasing strategy is similar to the epsilon greed but the epsilon value
decreases over time. This implementation provides a decreasing e0 / t  where e0 is a positive tuning
parameter and t the current round index.
The epsilon decreasing strategy is analysed in <i>Finite time analysis of the multiarmed
bandit problem.</i> by Auer, Cesa-Bianchi and Fisher in <i>Machine Learning</i> (2002).
```go
variants := []Variant{bandit.NewVariant(), bandit.NewVariant(), bandit.NewVariant()}
v, _ := bandit.EpsilonDecreasing(1000,variants)
```
* *EpsilonGreedy*:
This is the most simple of strategies. Intuitively it basically always pulls the lever with
the highest estimated mean, except when a randome lever is pulled with an epsilon frequency.
The function will return nil if len(variants) == 0
```go
variants := []Variant{bandit.NewVariant(), bandit.NewVariant(), bandit.NewVariant()}
v, _ := bandit.EpsilonGreedy(0.1, variants)
```

Docmentation
-----------
Documentation is available [on godoc](http://godoc.org/github.com/ghais/bandit).
