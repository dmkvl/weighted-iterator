package weightediterator

import (
	"math/rand"
)

type Iterator interface {
	HasNext() bool
	Next() int
	GetWeight() int
}

type iteratorFromSequence struct {
	sequence []int
	offset   int
	weight   int
}

func (it *iteratorFromSequence) HasNext() bool {
	return it.offset < len(it.sequence)
}

func (it *iteratorFromSequence) Next() int {
	next := it.sequence[it.offset]
	it.offset++
	return next
}

func (it *iteratorFromSequence) GetWeight() int {
	return it.weight
}

type WeightedIterator struct {
	totalWeight int
	// zero-index slice
	iterators []Iterator
	// zero-index slice, values are corresponding accumulated weights of iterators
	accumulatedWeights []int
}

// NewWeightedIterator returns the weighted iterator of iterators.
func NewWeightedIterator(config []IteratorConfig) *WeightedIterator {

	iterator := &WeightedIterator{
		iterators:          make([]Iterator, 0, len(config)),
		accumulatedWeights: make([]int, 0, len(config)),
	}

	for _, item := range config {
		iterator.totalWeight += item.Weight

		// Initializing iterators from config
		// (act as a primitive source of entropy).
		// There may probably be other iterator implementations.
		iterator.iterators = append(iterator.iterators, &iteratorFromSequence{
			sequence: item.Sequence,
			weight:   item.Weight,
		})

		// calculating the iterators accumulated weights
		iterator.accumulatedWeights = append(iterator.accumulatedWeights, iterator.totalWeight)
	}

	return iterator
}

// HasNext returns true if the next element exists.
// HasNext is not thread-safe.
func (rg *WeightedIterator) HasNext() bool {
	return len(rg.iterators) > 0 && rg.totalWeight > 0
}

// Next returns the next element.
// Next is not thread-safe.
func (rg *WeightedIterator) Next() int {
	// choose an iterator pseudo-randomly taking into account the weight
	randWeight := rand.Intn(rg.totalWeight)
	iteratorIndex := getNextRightAccumulatedWeightIndex(0, len(rg.accumulatedWeights)-1, randWeight, rg.accumulatedWeights)
	iter := rg.iterators[iteratorIndex]

	next := iter.Next()
	if !iter.HasNext() { // the iterator is exhausted
		rg.removeExhaustedIterator(iteratorIndex)
	}
	return next
}

// removeExhaustedIterator removes the exhausted iterator
// and recalculates the accumulated weights
func (rg *WeightedIterator) removeExhaustedIterator(iteratorIndex int) {
	if len(rg.iterators) == 1 {
		rg.totalWeight = 0
		rg.iterators = rg.iterators[:0]
		return
	}

	// we should not care about the iterators order here,
	// since the accumulatedWeights will be recalculated
	lastIndex := len(rg.iterators) - 1
	rg.iterators[iteratorIndex] = rg.iterators[lastIndex]
	rg.iterators[lastIndex] = nil
	rg.iterators = rg.iterators[:lastIndex]
	rg.accumulatedWeights = rg.accumulatedWeights[:lastIndex]

	rg.recalculateAccumulatedWeight()
}

// Each time an exhausted iterator is deleted, the accumulatedWeights must be recalculated.
func (rg *WeightedIterator) recalculateAccumulatedWeight() {
	rg.totalWeight = 0

	for index, iter := range rg.iterators {
		rg.totalWeight += iter.GetWeight()
		rg.accumulatedWeights[index] = rg.totalWeight
	}
}

// returns the index of the closest weight value to the right from the accumulatedWeights.
// For randomWeight it finds the next existing accumulated weight (using slightly modified binary search),
// in the accumulatedWeights slice and returns its index.
// e.g.:
// [0,  1,  2,  3,  4,  5] <-- accumulatedWeight indices
// [7, 14, 23, 40, 86, 91] <-- accumulatedWeights values
//
// for randomWeight == 23 it returns 2
// for randomWeight == 24 it returns 3
// explanation: [.., 23, 24 ---> 40, ...]
func getNextRightAccumulatedWeightIndex(start, end int, randomWeight int, accumulatedWeights []int) int {
	var middle int

	for end-start > 1 {
		middle = (end-start)/2 + start

		if accumulatedWeights[middle] == randomWeight {
			return middle
		}

		if accumulatedWeights[middle] < randomWeight {
			start = middle
		} else {
			end = middle
		}
	}

	if accumulatedWeights[start] >= randomWeight {
		return start
	}
	return end
}
