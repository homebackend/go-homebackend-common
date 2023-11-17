package homecommon

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

func GetConf[C any](configFile string) *C {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("yaml file read err   #%v ", err)
	}

	var c C
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	validate := validator.New()
	err = validate.Struct(c)
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	return &c
}
