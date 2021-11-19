package model

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

func (p *PublisherList) GetNames() []string {
	var v []string
	for _, i := range p.PubList {
		v = append(v, i.Name)
	}
	return v
}
