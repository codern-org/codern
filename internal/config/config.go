package config

import (
	"fmt"
	"os"

	"github.com/codern-org/codern/domain"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

func Load(path string) (*domain.Config, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config *domain.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validatePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory, not a file", path)
	}

	return nil
}
