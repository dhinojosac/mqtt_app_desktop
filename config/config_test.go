package config

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/dhinojosac/mqtt_app_desktop/model"
)

func TestReadBrokers(t *testing.T) {
	file, _ := os.Open("../conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	brokerList := []model.Broker{}
	err := decoder.Decode(&brokerList)
	if err != nil {
		t.Errorf("Error at decoding")
	}
	t.Logf("%v", brokerList)
}

func TestReadPublisher(t *testing.T) {
	file, _ := os.Open("../conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	brokerList := []model.Broker{}
	err := decoder.Decode(&brokerList)
	if err != nil {
		t.Errorf("Error at decoding")
	}
	publisher := brokerList[0].Publisher
	t.Logf("Publisher: %v", publisher)
}

func TestReadPubsFromBroker(t *testing.T) {
	s := "Demo Broker TGC"
	file, _ := os.Open("../conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	brokerList := []model.Broker{}
	err := decoder.Decode(&brokerList)
	if err != nil {
		t.Logf("Error at decoding")
	}

	for _, i := range brokerList {
		t.Logf("%v", i.Name)
		if strings.Contains(i.Name, s) {
			t.Logf("Publisher: %v", i.Publisher)
			return
		}
	}
	t.Errorf("Broker not found")

}
