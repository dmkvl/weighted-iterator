package weightediterator

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type WeightedIteratorSuite struct {
	suite.Suite

	iterator       *WeightedIterator
	configIntsList []int
}

func TestWeightedIteratorSuite(t *testing.T) {
	suite.Run(t, new(WeightedIteratorSuite))
}

func (s *WeightedIteratorSuite) SetupTest() {

	config, err := LoadFromFile("config/example_iterators.yml")

	s.NoError(err)
	s.NoError(config.Validate())

	var totalLen int
	for _, c := range config.Iterators {
		totalLen += len(c.Sequence)
	}

	configInts := make([]int, 0, totalLen)
	for _, c := range config.Iterators {
		configInts = append(configInts, c.Sequence...)
	}

	s.iterator = NewWeightedIterator(config.Iterators)
	s.configIntsList = configInts
}

func (s *WeightedIteratorSuite) TestConvergence() {

	result := make([]int, 0, len(s.configIntsList))

	for s.iterator.HasNext() {
		result = append(result, s.iterator.Next())
	}

	s.ElementsMatchf(result, s.configIntsList, "lists are not equal: %v, %v", result, s.configIntsList)
}

func (s *WeightedIteratorSuite) TestGetNextRightAccumulatedWeightIndex() {

	var suits = []struct {
		name               string
		randWeight         int
		accumulatedWeights []int
		expectedIndex      int
	}{
		{
			"first weight value of the leftmost segment",
			1,
			[]int{7, 14, 23, 40, 86, 91},
			0,
		},
		{
			"penultimate weight value of the leftmost segment",
			6,
			[]int{7, 14, 23, 40, 86, 91},
			0,
		},
		{
			"the last weight value of the leftmost segment",
			7,
			[]int{7, 14, 23, 40, 86, 91},
			0,
		},
		{
			"the first weight value of the second segment",
			8,
			[]int{7, 14, 23, 40, 86, 91},
			1,
		},
		{
			"the first weight value of the most right segment",
			87,
			[]int{7, 14, 23, 40, 86, 91},
			5,
		},
		{
			"the last weight value of the most right segment",
			91,
			[]int{7, 14, 23, 40, 86, 91},
			5,
		},
		{
			"the first weight value of the single weight segment",
			0,
			[]int{5},
			0,
		},
		{
			"the last weight value of the single weight segment",
			5,
			[]int{5},
			0,
		},
	}

	for _, suit := range suits {
		s.Suite.Run(suit.name, func() {
			index := getNextRightAccumulatedWeightIndex(0, len(suit.accumulatedWeights)-1, suit.randWeight, suit.accumulatedWeights)
			s.Equal(suit.expectedIndex, index)
		})
	}
}
