package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dhinojosac/mqtt_app_desktop/model"
)

// func ReadPublisher() []model.PublisherData {
// 	file, _ := os.Open("conf.json")
// 	defer file.Close()

// 	decoder := json.NewDecoder(file)
// 	pubList := []model.PublisherData{}
// 	err := decoder.Decode(&pubList)
// 	if err != nil {
// 		fmt.Println("error:", err)
// 	}
// 	fmt.Printf("%v", pubList)
// 	return pubList

// }

// func WritePublisher(m []model.PublisherData) {
// 	file, _ := json.MarshalIndent(m, "", " ")

// 	_ = ioutil.WriteFile("conf.json", file, 0644)

// }

func ReadBrokers() []model.Broker {
	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	brokerList := []model.Broker{}
	err := decoder.Decode(&brokerList)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Printf("%v", brokerList)
	return brokerList
}

func WriteBrokers(m []model.Broker) {
	file, _ := json.MarshalIndent(m, "", " ")

	_ = ioutil.WriteFile("conf.json", file, 0644)

}

// Store new broker in the list and replace if it exists
func WriteBroker(m model.Broker) {
	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	brokerList := []model.Broker{}
	err := decoder.Decode(&brokerList)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%v", brokerList)
	for i, b := range brokerList {
		if b.Name == m.Name {
			brokerList[i] = m
			break
		}
	}
	WriteBrokers(brokerList)
}

// Delete broker from the list
func DeleteBroker(s string) {
	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	brokerList := []model.Broker{}
	err := decoder.Decode(&brokerList)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%v", brokerList)
	for i, b := range brokerList {
		if b.Name == s {
			brokerList = append(brokerList[:i], brokerList[i+1:]...)
			break
		}
	}
	WriteBrokers(brokerList)
}

func ReadBroker(s string) model.Broker {
	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	brokerList := []model.Broker{}
	err := decoder.Decode(&brokerList)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%v", brokerList)
	for _, i := range brokerList {
		if s == i.Name {
			return i
		}
	}
	return model.Broker{}
}

func ReadPubsFromBroker(s string) []model.PublisherData {
	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	brokerList := []model.Broker{}
	err := decoder.Decode(&brokerList)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%v", brokerList)
	for _, i := range brokerList {
		if s == i.Name {
			return i.Publisher
		}
	}
	return []model.PublisherData{}
}
