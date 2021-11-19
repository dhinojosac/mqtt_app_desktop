// Package main provides various examples of Fyne API capabilities
package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dhinojosac/mqtt_app_desktop/config"
	"github.com/dhinojosac/mqtt_app_desktop/model"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var isConnected bool
var isSubscribed bool
var Client MQTT.Client

var containerList *fyne.Container
var desub_button *widget.Button
var subs_button *widget.Button
var checkAuto *widget.Check

var pubList model.PublisherList

var autoPublish bool

func checkConnection() {
	for {
		//fmt.Printf("Client status: %v\n", Client.IsConnected())
		time.Sleep(5 * time.Second)
		if !Client.IsConnected() {
			Client.Connect()
			break
		}
	}
}

func MQTTSetting(win fyne.Window) fyne.CanvasObject {
	mqtt_address := widget.NewEntry()
	mqtt_address.SetPlaceHolder("cloud.thegrouplab.com")
	mqtt_address.Text = "cloud.thegrouplab.com"
	mqtt_user := widget.NewEntry()
	mqtt_user.SetPlaceHolder("tgl")
	mqtt_user.SetText("tgl")
	mqtt_password := widget.NewPasswordEntry()
	mqtt_password.SetPlaceHolder("tgl1234id")
	mqtt_password.SetText("tgl1234id")

	mqtt_log := widget.NewLabel("")
	var conn_button *widget.Button
	conn_button = widget.NewButton("Connect", func() {
		if !isConnected {
			if mqtt_address.Text == "" {
				err := errors.New("MQTT Broker empty")
				dialog.ShowError(err, win)
			} else {
				opts := MQTT.NewClientOptions().AddBroker("tcp://" + mqtt_address.Text + ":1883") //"tcp://code4fun.cl:1883"
				if mqtt_user.Text != "" {
					opts.SetUsername(mqtt_user.Text)
				}
				if mqtt_password.Text != "" {
					opts.SetPassword(mqtt_password.Text)
				}
				Client = MQTT.NewClient(opts)
				if token := Client.Connect(); token.Wait() && token.Error() != nil {
					log.Fatal(token.Error())
				}
				isConnected = true
				mqtt_log.SetText("Client connected")
				conn_button.SetText("Disconnect")
				mqtt_address.Disable()
				mqtt_user.Disable()
				mqtt_password.Disable()
				checkAuto.Enable()
				go checkConnection()
			}
		} else {
			Client.Disconnect(200)
			isConnected = false
			mqtt_log.SetText("Client disconnected")
			conn_button.SetText("Connect")
			mqtt_address.Enable()
			mqtt_user.Enable()
			mqtt_password.Enable()
		}

	})
	return container.New(layout.NewVBoxLayout(),
		widget.NewLabel("MQTT Broker"),
		mqtt_address,
		widget.NewLabel("User"),
		mqtt_user,
		widget.NewLabel("Password"),
		mqtt_password,
		conn_button,
		mqtt_log,
	)
}

func MQTTPubScreen(win fyne.Window) fyne.CanvasObject {

	mqtt_topic := widget.NewEntry()
	mqtt_topic.SetPlaceHolder("devices/<organization>/device/<identifier>/measurements")
	mqtt_topic.Text = "devices/2/device/DEMO2/measurements"

	mqtt_message_label := widget.NewLabel("Message:")
	mqtt_message_input := widget.NewMultiLineEntry()
	mqtt_message_input.SetText(`{"temp":24, "hum":50}`)

	mqtt_log := widget.NewLabel("")
	pub_button := widget.NewButton("Publish", func() {
		if isConnected {
			if mqtt_topic.Text == "" {
				err := errors.New("MQTT Topic empty")
				dialog.ShowError(err, win)
			} else {
				fmt.Println("Publish message")

				if token := Client.Publish(mqtt_topic.Text, 0, false, mqtt_message_input.Text); token.Wait() && token.Error() != nil {
					log.Fatal(token.Error())
				}
				mqtt_log.SetText("Message sent!")
			}
		} else {
			fmt.Println("Client not connected")
			mqtt_log.SetText("Client not connected")
		}

	})
	autoTime := widget.NewEntry()
	autoTime.SetText("10")
	var ticker *time.Ticker

	checkAuto = widget.NewCheck("Automate", func(val bool) {
		autoPublish = val
		if autoPublish {
			fmt.Println("Check pressed")
			pub_button.Disable()
			mqtt_message_input.Disable()
			mqtt_topic.Disable()

			autoTimeInt, err := strconv.Atoi(autoTime.Text)
			if err != nil {
				fmt.Printf("Error")
			}

			go func() {
				ticker = time.NewTicker(time.Duration(autoTimeInt) * time.Second)
				for range ticker.C {
					if token := Client.Publish(mqtt_topic.Text, 0, false, mqtt_message_input.Text); token.Wait() && token.Error() != nil {
						log.Fatal(token.Error())
					}
				}
			}()

		} else {

			ticker.Stop()
			pub_button.Enable()
			mqtt_message_input.Enable()
			mqtt_topic.Enable()
		}

	})
	checkAuto.Disable()

	pub_adv := container.New(layout.NewHBoxLayout(), widget.NewLabel("Time:"), autoTime, widget.NewLabel("s"), widget.NewLabel("  "), checkAuto)

	// TODO: get publisher from file
	// v1 := model.PublisherData{Name: "Demo2", Topic: "devices/2/device/DEMO2/measurements", Message: `{"temp":30, "hum":55, "bat":72}`}
	// v2 := model.PublisherData{Name: "Vitalis Test", Topic: "topic/2/2", Message: `{"temp":24, "spo2": 50, "battery_level":70, "sbp":70,"dbp":110, "respiration_rate":18,"heart_rate":72}`}
	// v3 := model.PublisherData{Name: "Custom", Topic: "test/alert", Message: `{"alert":0}`}
	// pubList.SetItems([]model.PublisherData{v1, v2, v3})

	pubList.SetItems(config.ReadPublisher()) //Read conf.json

	publisherSelect := widget.NewSelectEntry(pubList.GetNames())
	publisherSelect.OnChanged = func(s string) {
		p := pubList.GetItem(s)
		mqtt_topic.SetText(p.Topic)
		mqtt_message_input.SetText(p.Message)
	}
	// end publisher implementation

	return container.New(layout.NewVBoxLayout(),
		widget.NewLabel("Topic:"),
		mqtt_topic,
		mqtt_message_label,
		mqtt_message_input,
		pub_button,
		widget.NewLabel("Select Publish Template:"),
		publisherSelect,
		pub_adv,
		mqtt_log,
	)
}

func MQTTSubScreen(win fyne.Window) fyne.CanvasObject {
	mqtt_topic := widget.NewEntry()
	mqtt_topic.SetPlaceHolder("#")
	mqtt_topic.Text = "#"
	//mqtt_message_label := widget.NewLabel("Incoming Messages:")

	//scroll

	containerList = container.New(layout.NewVBoxLayout(), widget.NewLabel(""))
	containerList.Resize(fyne.NewSize(800, 100))
	index := 1

	scroll := container.NewVScroll(containerList)
	scroll.Resize(fyne.NewSize(800, 100))

	desub_button = widget.NewButton("Unsubscribe", func() {
		if isConnected && isSubscribed {
			Client.Unsubscribe(mqtt_topic.Text)
			isSubscribed = false
			subs_button.Enable()
			mqtt_topic.Enable()
			containerList.Add(widget.NewLabel("[!] Unsubscribed to: " + mqtt_topic.Text))
			scroll.ScrollToBottom()

		} else {
			fmt.Println("Client not connected")
			containerList.Add(widget.NewLabel("Client not connected"))
			scroll.ScrollToBottom()
		}
		if !isSubscribed {
			desub_button.Disable()
		}

	})

	desub_button.Disable()

	subs_button = widget.NewButton("Subscribe", func() {
		if isConnected {
			if mqtt_topic.Text == "" {
				err := errors.New("[!] MQTT Topic empty")
				dialog.ShowError(err, win)
			} else {
				fmt.Println("Subscribe")
				subs_button.Text = "Subscribed"
				subs_button.Disable()
				mqtt_topic.Disable()
				containerList.Add(widget.NewLabel("[!] Subscribed to: " + mqtt_topic.Text))
				scroll.ScrollToBottom()
				go func() {
					isSubscribed = true
					desub_button.Enable()
					Client.Subscribe(mqtt_topic.Text, 0, func(client MQTT.Client, msg MQTT.Message) {
						console_text := fmt.Sprintf("[%d] topic: %s\nmsg:%s", index, string(msg.Topic()), string(msg.Payload()))
						containerList.Add(widget.NewLabel(console_text))
						scroll.ScrollToBottom()
						index++
					})
				}()

			}

		} else {
			fmt.Println("Client not connected")
			containerList.Add(widget.NewLabel("Client not connected"))
		}

	})
	v1 := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("Topic"),
		mqtt_topic,
		subs_button,
		desub_button,
	)

	return container.New(layout.NewGridLayoutWithRows(2), v1, scroll)
}

func makeUI() fyne.CanvasObject {
	f := binding.NewFloat()
	prog := widget.NewProgressBarWithData(f)
	slide := widget.NewSliderWithData(0, 1, f)
	go func(f binding.Float) {
		for {
			f1, _ := f.Get()
			f.Set(f1 + 0.1)
			time.Sleep(1 * time.Second)
		}
	}(f)

	slide.Step = 0.01
	btn := widget.NewButton("Set to 0.0", func() {
		_ = f.Set(0.0)
	})

	a1 := container.NewVBox(prog, slide)
	a2 := container.New(layout.NewGridLayoutWithRows(2), a1, btn)

	return a2

}

func main() {
	//gui
	a := app.New()
	w := a.NewWindow("TGC MQTT Tester v0.1")

	w.SetIcon(appLogo)

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())

	w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
		fyne.NewMenuItem("New", func() { fmt.Println("Menu New") }),
		// a quit item will be appended to our first menu
	), fyne.NewMenu("Edit",
		fyne.NewMenuItem("Cut", func() { fmt.Println("Menu Cut") }),
		fyne.NewMenuItem("Copy", func() { fmt.Println("Menu Copy") }),
		fyne.NewMenuItem("Paste", func() { fmt.Println("Menu Paste") }),
	)))

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Config", theme.SettingsIcon(), MQTTSetting(w)),
		container.NewTabItemWithIcon("Publish", theme.MailSendIcon(), MQTTPubScreen(w)),
		container.NewTabItemWithIcon("Subscribe", theme.SearchIcon(), MQTTSubScreen(w)),
		//container.NewTabItemWithIcon("Test", theme.ComputerIcon(), makeUI()),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(900, 320))
	w.ShowAndRun()
}
