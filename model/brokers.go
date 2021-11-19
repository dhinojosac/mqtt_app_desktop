package model

type Broker struct {
	Name        string
	Description string
	Endpoint    string
	User        string
	Password    string
}

type BrokerList []Broker
