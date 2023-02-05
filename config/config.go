package config

import (
	"aviso/domain"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func GetTargets(pathconfig string) (*domain.Targets, error) {
	data, err := os.ReadFile(pathconfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config yaml")
	}
	var t domain.Targets
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		log.Println("Unmarshalling error: ")
		return nil, err
	}

	return &t, nil
}
