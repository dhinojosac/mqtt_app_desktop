package model

import "fmt"

type Broker struct {
	Name        string
	Description string
	Endpoint    string
	User        string
	Password    string
	Port        string
	Publisher   []PublisherData `json:"pubs"`
}

type BrokerList []Broker

func (b *BrokerList) GetItem(s string) Broker {
	for _, i := range *b {
		if s == i.Name {
			return i
		}
	}
	return Broker{}
}

func (b *BrokerList) SetItems(s []Broker) {
	*b = append(*b, s...)
}

func (b *BrokerList) ReSetItems(s []Broker) {
	*b = nil
	*b = append(*b, s...)
}

func (b *BrokerList) GetNames() []string {
	var v []string
	for _, i := range *b {
		v = append(v, i.Name)
	}
	return v
}

func (b *BrokerList) AddItem(s Broker) {
	for _, i := range *b {
		if i.Name == s.Name {
			s.Name = fmt.Sprintf("%s %s", s.Name, "(Copy)")
		}
	}
	*b = append(*b, s)
}

func (b *BrokerList) DeleteItem(s string) {

	for i, n := range *b {
		if n.Name == s {
			*b = append((*b)[:i], (*b)[i+1:]...)
		}
	}
}

func (b *Broker) GetPublisherNames() []string {
	var v []string
	for _, i := range b.Publisher {
		v = append(v, i.Name)
	}
	fmt.Printf("GetPublisherNames: %v", v)
	return v
}

func (b *Broker) GetPublisher(s string) PublisherData {
	for _, i := range b.Publisher {
		if s == i.Name {
			return i
		}
	}
	return PublisherData{}
}

func (b *Broker) AddPublisher(s PublisherData) {
	for _, i := range b.Publisher {
		if i.Name == s.Name {
			s.Name = fmt.Sprintf("%s %s", s.Name, "(Copy)")
		}
	}
	b.Publisher = append(b.Publisher, s)
}

func (b *Broker) DeletePublisher(s string) {

	for i, n := range b.Publisher {
		if n.Name == s {
			b.Publisher = append(b.Publisher[:i], b.Publisher[i+1:]...)
		}
	}
}
