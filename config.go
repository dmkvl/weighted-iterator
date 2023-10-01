package weightediterator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

type IteratorConfig struct {
	Weight   int
	Sequence []int
}

func (ic IteratorConfig) Validate() error {
	return validation.ValidateStruct(&ic,
		validation.Field(&ic.Weight, validation.Required, validation.Min(1)),
		validation.Field(&ic.Sequence, validation.Required, validation.Length(1, 0)),
	)
}

type Config struct {
	Iterators []IteratorConfig `konaf:"iterators"`
}

func (c Config) Validate() error {
	for _, item := range c.Iterators {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func LoadFromFile(filePath string) (Config, error) {
	return loadConfig(file.Provider(filePath))
}

func LoadFromBytes(rawConfig []byte) (Config, error) {
	return loadConfig(rawbytes.Provider(rawConfig))
}

func loadConfig(provider koanf.Provider) (Config, error) {
	var k = koanf.NewWithConf(koanf.Conf{
		Delim:       ".",
		StrictMerge: true,
	})

	if err := k.Load(provider, yaml.Parser()); err != nil {
		return Config{}, err
	}

	var config Config

	if err := k.Unmarshal("", &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
