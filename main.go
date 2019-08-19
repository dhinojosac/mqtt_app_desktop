// Package main provides various examples of Fyne API capabilities
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MQTTScreen(win fyne.Window) fyne.CanvasObject {

	mqtt_address := widget.NewEntry()
	mqtt_address.SetPlaceHolder("code4fun.cl")
	mqtt_user := widget.NewEntry()
	mqtt_user.SetPlaceHolder("JohnDoe")
	mqtt_password := widget.NewPasswordEntry()
	mqtt_password.SetPlaceHolder("*****")
	mqtt_topic := widget.NewEntry()
	mqtt_topic.SetPlaceHolder("test")
	mqtt_message_label := widget.NewLabel("Message:")
	mqtt_message_input := widget.NewMultiLineEntry()
	send_button := widget.NewButton(fmt.Sprintf("Publish"), func() {
		if mqtt_topic.Text == "" {
			err := errors.New("MQTT Topic empty!")
			dialog.ShowError(err, win)
		} else if mqtt_address.Text == "" {
			err := errors.New("MQTT Broker empty!")
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
		send_button,
	)
}

func MQTTSubScreen(win fyne.Window) fyne.CanvasObject {
	mqtt_address := widget.NewEntry()
	mqtt_address.SetPlaceHolder("code4fun.cl")
	mqtt_user := widget.NewEntry()
	mqtt_user.SetPlaceHolder("JohnDoe")
	mqtt_password := widget.NewPasswordEntry()
	mqtt_password.SetPlaceHolder("*****")
	mqtt_topic := widget.NewEntry()
	mqtt_topic.SetPlaceHolder("test")
	mqtt_message_label := widget.NewLabel("Messages:")
	list := widget.NewVBox()
	scroll := widget.NewScrollContainer(list)
	scroll.Resize(fyne.NewSize(200, 200))
	index := 1
	send_button := widget.NewButton(fmt.Sprintf("Subscribe"), func() {
		if mqtt_topic.Text == "" {
			err := errors.New("MQTT Topic empty!")
			dialog.ShowError(err, win)
		} else if mqtt_address.Text == "" {
			err := errors.New("MQTT Broker empty!")
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

			go func() {
				client.Subscribe(mqtt_topic.Text, 0, func(client mqtt.Client, msg mqtt.Message) {
					console_text := fmt.Sprintf(">topic:[%s] message: %s\n", msg.Topic(), string(msg.Payload()))
					list.Append(widget.NewLabel(console_text))
					index++
				})
			}()

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
		send_button,
		mqtt_message_label,
		scroll,
	)
}

func main() {
	//mqtt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

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
	w.Resize(fyne.NewSize(520, 160))
	w.SetFixedSize(true)
	w.ShowAndRun()
}
