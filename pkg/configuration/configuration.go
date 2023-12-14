package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func NewConfig[TConfig any](path, configName, environment string) (*TConfig, error) {
	formattedPath := fmt.Sprintf("%s/%s.%s.json", path, configName, environment)
	file, err := os.Open(formattedPath)
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	decoder := json.NewDecoder(file)
	configuration := new(TConfig)
	err = decoder.Decode(configuration)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return configuration, nil
}
