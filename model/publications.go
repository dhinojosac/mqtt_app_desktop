package model

import "fmt"

type PublisherData struct {
	Name        string
	Description string
	Topic       string
	Message     string
}

type PublisherList struct {
	PubList []PublisherData
}

func (p *PublisherList) GetItem(s string) PublisherData {
	for _, i := range p.PubList {
		if s == i.Name {
			return i
		}
	}
	return PublisherData{}
}

func (p *PublisherList) SetItems(s []PublisherData) {
	p.PubList = append(p.PubList, s...)
}

func (p *PublisherList) ReSetItems(s []PublisherData) {
	p.PubList = nil
	p.PubList = append(p.PubList, s...)
}

func (p *PublisherList) GetNames() []string {
	var v []string
	for _, i := range p.PubList {
		v = append(v, i.Name)
	}
	return v
}

func (p *PublisherList) AddItem(s PublisherData) {
	for _, i := range p.PubList {
		if i.Name == s.Name {
			s.Name = fmt.Sprintf("%s %s", s.Name, "(Copy)")
		}
	}
	p.PubList = append(p.PubList, s)
}

func (p *PublisherList) DeleteItem(s string) {

	for i, n := range p.PubList {
		if n.Name == s {
			p.PubList = append(p.PubList[:i], p.PubList[i+1:]...)
		}
	}
}
