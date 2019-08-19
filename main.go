// Package main provides various examples of Fyne API capabilities
package main

import (
	"errors"
	"fmt"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func welcomeScreen(a fyne.App) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(theme.FyneLogo())
	logo.SetMinSize(fyne.NewSize(320, 320))
	link, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return widget.NewVBox(
		widget.NewLabelWithStyle("Welcome to the Fyne toolkit demo app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewHyperlinkWithStyle("fyne.io", link, fyne.TextAlignCenter, fyne.TextStyle{}),
		layout.NewSpacer(),

		widget.NewGroup("Theme",
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewButton("Dark", func() {
					a.Settings().SetTheme(theme.DarkTheme())
				}),
				widget.NewButton("Light", func() {
					a.Settings().SetTheme(theme.LightTheme())
				}),
			),
		),
	)
}

func MQTTScreen(win fyne.Window) fyne.CanvasObject {
	mqtt_address := widget.NewEntry()
	mqtt_address.SetPlaceHolder("MQTT Broker Address")
	mqtt_password := widget.NewPasswordEntry()
	mqtt_password.SetPlaceHolder("MQTT password")
	mqtt_topic := widget.NewEntry()
	mqtt_topic.SetPlaceHolder("MQTT Topic")
	mqtt_message_label := widget.NewLabel("Message:")
	mqtt_message_input := widget.NewMultiLineEntry()
	send_button := widget.NewButton(fmt.Sprintf("Publish"), func() {
		fmt.Println("Publish")
	})
	return widget.NewVBox(
		mqtt_address,
		mqtt_password,
		mqtt_topic,
		mqtt_message_label,
		mqtt_message_input,
		send_button,
	)
}

func MQTTSubScreen(win fyne.Window) fyne.CanvasObject {
	mqtt_address := widget.NewEntry()
	mqtt_address.SetPlaceHolder("MQTT Broker Address")
	mqtt_password := widget.NewPasswordEntry()
	mqtt_password.SetPlaceHolder("MQTT password")
	mqtt_topic := widget.NewEntry()
	mqtt_topic.SetPlaceHolder("MQTT Topic")
	mqtt_message_label := widget.NewLabel("Messages:")
	send_button := widget.NewButton(fmt.Sprintf("Subscribe"), func() {
		fmt.Println("Subscribe")
		if mqtt_topic.Text == "" {
			err := errors.New("MQTT Topic empty!")
			dialog.ShowError(err, win)
		}
		if mqtt_address.Text == "" {
			err := errors.New("MQTT Broker empty!")
			dialog.ShowError(err, win)
		}

	})
	return widget.NewVBox(
		mqtt_address,
		mqtt_password,
		mqtt_topic,
		send_button,
		mqtt_message_label,
	)
}

func main() {
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
		widget.NewTabItemWithIcon("MQTT Publish", theme.ContentCopyIcon(), MQTTScreen(w)),
		widget.NewTabItemWithIcon("MQTT Subscribe", theme.ContentCopyIcon(), MQTTSubScreen(w)),
	)
	tabs.SetTabLocation(widget.TabLocationLeading)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(480, 160))
	w.SetFixedSize(true)
	w.ShowAndRun()
}
