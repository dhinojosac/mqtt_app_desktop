// Package main provides various examples of Fyne API capabilities
package main

import (
	"errors"
	"fmt"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func MQTTScreen(win fyne.Window) fyne.CanvasObject {

	mqtt_address := widget.NewEntry()
	mqtt_address.SetPlaceHolder("code4fun.cl")
	mqtt_address.Text = "code4fun.cl"
	mqtt_user := widget.NewEntry()
	mqtt_user.SetPlaceHolder("JohnDoe")
	mqtt_password := widget.NewPasswordEntry()
	mqtt_password.SetPlaceHolder("*****")
	mqtt_topic := widget.NewEntry()
	mqtt_topic.SetPlaceHolder("test")
	mqtt_topic.Text = "test"
	mqtt_message_label := widget.NewLabel("Message:")
	mqtt_message_input := widget.NewMultiLineEntry()
	pub_button := widget.NewButton(fmt.Sprintf("Publish"), func() {
		if mqtt_address.Text == "" {
			err := errors.New("MQTT Broker empty!")
			dialog.ShowError(err, win)
		} else if mqtt_topic.Text == "" {
			err := errors.New("MQTT Topic empty!")
			dialog.ShowError(err, win)
		} else {
			fmt.Println("Publish message")
			opts := MQTT.NewClientOptions().AddBroker("tcp://" + mqtt_address.Text + ":1883") //"tcp://code4fun.cl:1883"
			if mqtt_user.Text != "" {
				opts.SetUsername(mqtt_user.Text)
			}
			if mqtt_password.Text != "" {
				opts.SetPassword(mqtt_password.Text)
			}
			client := MQTT.NewClient(opts)
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				log.Fatal(token.Error())
			}
			if token := client.Publish(mqtt_topic.Text, 0, false, mqtt_message_input.Text); token.Wait() && token.Error() != nil {
				log.Fatal(token.Error())
			}
		}
	})
	return widget.NewVBox(
		widget.NewLabel("MQTT Broker"),
		mqtt_address,
		widget.NewLabel("User"),
		mqtt_user,
		widget.NewLabel("Password"),
		mqtt_password,
		widget.NewLabel("Topic"),
		mqtt_topic,
		mqtt_message_label,
		mqtt_message_input,
		pub_button,
	)
}

func MQTTSubScreen(win fyne.Window) fyne.CanvasObject {
	mqtt_address := widget.NewEntry()
	mqtt_address.SetPlaceHolder("code4fun.cl")
	mqtt_address.Text = "code4fun.cl"
	mqtt_user := widget.NewEntry()
	mqtt_user.SetPlaceHolder("JohnDoe")
	mqtt_password := widget.NewPasswordEntry()
	mqtt_password.SetPlaceHolder("*****")
	mqtt_topic := widget.NewEntry()
	mqtt_topic.SetPlaceHolder("test")
	mqtt_topic.Text = "test"
	//mqtt_message_label := widget.NewLabel("Incoming Messages:")

	//scroll
	list := widget.NewVBox()
	list.Resize(fyne.NewSize(100, 100))
	index := 1
	/*
		//Fill scrollable list
		for i := 1; i <= 20; i++ {
			list.Append(widget.NewLabel(fmt.Sprintf("test %d", index)))
			index++
		}
	*/
	scroll := widget.NewScrollContainer(list)

	//button subscribe
	var subs_button *widget.Button
	subs_button = widget.NewButton(fmt.Sprintf("Subscribe"), func() {
		if mqtt_address.Text == "" {
			err := errors.New("MQTT Broker empty!")
			dialog.ShowError(err, win)
		} else if mqtt_topic.Text == "" {
			err := errors.New("MQTT Topic empty!")
			dialog.ShowError(err, win)
		} else {
			fmt.Println("Subscribe")
			opts := MQTT.NewClientOptions().AddBroker("tcp://" + mqtt_address.Text + ":1883") //"tcp://code4fun.cl:1883"
			if mqtt_user.Text != "" {
				opts.SetUsername(mqtt_user.Text)
			}
			if mqtt_password.Text != "" {
				opts.SetPassword(mqtt_password.Text)
			}
			client := MQTT.NewClient(opts)
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				log.Fatal(token.Error())
			}

			subs_button.Text = "Subscribed"
			subs_button.Disable()
			list.Append(widget.NewLabel("Subscribed to " + mqtt_address.Text))
			go func() {
				client.Subscribe(mqtt_topic.Text, 0, func(client MQTT.Client, msg MQTT.Message) {
					console_text := fmt.Sprintf("> %s", string(msg.Payload()))
					list.Append(widget.NewLabel(console_text))
					index++
				})
			}()

		}

	})
	v1 := widget.NewVBox(
		widget.NewLabel("MQTT Broker"),
		mqtt_address,
		widget.NewLabel("User"),
		mqtt_user,
		widget.NewLabel("Password"),
		mqtt_password,
		widget.NewLabel("Topic"),
		mqtt_topic,
		subs_button,
	)

	content := fyne.NewContainerWithLayout(layout.NewGridLayout(2), v1, scroll)

	return content
}

func main() {
	//gui
	a := app.New()
	w := a.NewWindow("MQTT App")
	a.Settings().SetTheme(theme.LightTheme())

	w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
		fyne.NewMenuItem("New", func() { fmt.Println("Menu New") }),
		// a quit item will be appended to our first menu
	), fyne.NewMenu("Edit",
		fyne.NewMenuItem("Cut", func() { fmt.Println("Menu Cut") }),
		fyne.NewMenuItem("Copy", func() { fmt.Println("Menu Copy") }),
		fyne.NewMenuItem("Paste", func() { fmt.Println("Menu Paste") }),
	)))

	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("MQTT Publish", theme.MoveUpIcon(), MQTTScreen(w)),
		widget.NewTabItemWithIcon("MQTT Subscribe", theme.MoveDownIcon(), MQTTSubScreen(w)),
	)

	tabs.SetTabLocation(widget.TabLocationLeading)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(720, 320))
	w.ShowAndRun()
}
