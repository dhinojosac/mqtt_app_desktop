package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ChatScrollHistory struct {
	Container *fyne.Container
	Scroll    *container.Scroll
	Messages  []string
}

func NewChatScrollHistory() *ChatScrollHistory {
	t := fmt.Sprintf("[%s]\n[!] Remember to Subscribe to a topic to start receiving messages\n", time.Now().String())
	v := container.New(layout.NewVBoxLayout(), widget.NewLabel(t))
	c := &ChatScrollHistory{
		Container: v,
		Scroll:    container.NewScroll(v),
		Messages:  []string{},
	}
	return c
}

func (c *ChatScrollHistory) AddMessage(msg string) {
	if len(c.Messages) >= 20 {
		//remove old message of slice
		c.Messages = append(c.Messages[1:], msg)
	}
	c.Messages = append(c.Messages, msg)
	c.Container.Add(widget.NewLabel(msg))
	c.Scroll.ScrollToBottom()
}

//Function that returns length of the messages array
func (c *ChatScrollHistory) GetLength() int {
	return len(c.Messages)
}

//Clear the messages array
func (c *ChatScrollHistory) Clear() {
	c.Messages = nil
	c.Container.Objects = nil
	c.Messages = []string{}
	c.Container.Objects = []fyne.CanvasObject{}
	c.Container.Refresh()
	c.Scroll.Refresh()
	c.Scroll.ScrollToTop()
}
