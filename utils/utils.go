package utils

import (
	"io/ioutil"
	"os"
)

func ReadFile(filepath string) ([]byte, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// var resp []string

	// err = json.Unmarshal(data, &resp)
	// if err != nil {
	// 	log.Fatalf("Failed to read file: %v", err)
	// }

	return data, nil
}
