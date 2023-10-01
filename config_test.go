package weightediterator

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) TestLoadConfig() {
	config, err := LoadFromFile("config/example_iterators.yml")

	s.NoError(err)
	s.NoError(config.Validate())
}
