package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dhinojosac/mqtt_app_desktop/model"
)

func ReadPublisher() []model.PublisherData {
	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	pubList := []model.PublisherData{}
	err := decoder.Decode(&pubList)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%v", pubList)
	return pubList

}

func WritePublisher(m []model.PublisherData) {
	file, _ := json.MarshalIndent(m, "", " ")

	_ = ioutil.WriteFile("conf.json", file, 0644)

}
