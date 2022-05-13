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
	"github.com/dhinojosac/mqtt_app_desktop/ui"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/tidwall/pretty"
)

var MAX_LINES = 200

var isConnected bool
var isSubscribed bool
var Client MQTT.Client

var containerList *fyne.Container //store history on subscribe view
var desub_button *widget.Button
var subs_button *widget.Button
var checkAuto *widget.Check
var publisherSelect *widget.SelectEntry

var pubItem *container.TabItem
var subItem *container.TabItem
var logItem *container.TabItem
var AppTabs *container.AppTabs

var brokerList model.BrokerList // list of brokers in memory
var currentBroker model.Broker  // current broker in memory

var autoPublish bool

var msgChan = make(chan MQTT.Message, 10)

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
	mqtt_address.SetPlaceHolder("Select your broker")
	mqtt_address.SetPlaceHolder("mqtt.broker.com")
	mqtt_user := widget.NewEntry()
	mqtt_user.SetPlaceHolder("username")
	mqtt_password := widget.NewPasswordEntry()
	mqtt_password.SetPlaceHolder("password")
	mqtt_port := widget.NewEntry()
	mqtt_port.SetPlaceHolder("1883")

	brokerList.SetItems(config.ReadBrokers()) //Read conf.json

	brokerSelect := widget.NewSelectEntry(brokerList.GetNames())
	brokerSelect.OnChanged = func(s string) {
		p := brokerList.GetItem(s)
		mqtt_address.SetText(p.Endpoint)
		mqtt_user.SetText(p.User)
		mqtt_password.SetText(p.Password)
		mqtt_port.SetText(p.Port)
		currentBroker = config.ReadBroker(s)
		fmt.Printf("\nSelected broker: %s\n", currentBroker.Endpoint)
		publisherSelect.SetOptions(currentBroker.GetPublisherNames())
	}

	mqtt_log := widget.NewLabel("")
	var conn_button *widget.Button
	conn_button = widget.NewButton("CONNECT", func() {
		if !isConnected {
			if mqtt_address.Text == "" {
				err := errors.New("MQTT Broker empty")
				dialog.ShowError(err, win)
			} else {
				opts := MQTT.NewClientOptions().AddBroker("tcp://" + mqtt_address.Text + ":" + mqtt_port.Text) //"tcp://code4fun.cl:1883"
				if mqtt_user.Text != "" {
					opts.SetUsername(mqtt_user.Text)
				}
				if mqtt_password.Text != "" {
					opts.SetPassword(mqtt_password.Text)
				}
				//generate client id
				clientName := "tgc-mqtt-tester-client" + strconv.Itoa(int(time.Now().Unix()))
				// add random client id
				opts.SetClientID(clientName)

				Client = MQTT.NewClient(opts)
				if token := Client.Connect(); token.Wait() && token.Error() != nil {
					log.Fatal(token.Error())
				}
				isConnected = true
				mqtt_log.SetText("Client connected to " + mqtt_address.Text + ":" + mqtt_port.Text + "\nwith client name: " + clientName)
				conn_button.SetText("DISCONNECT")

				mqtt_address.Disable()
				mqtt_user.Disable()
				mqtt_password.Disable()
				mqtt_port.Disable()
				checkAuto.Enable()
				go checkConnection()
			}
		} else {
			Client.Disconnect(200)
			isConnected = false
			mqtt_log.SetText("Client disconnected")
			conn_button.SetText("CONNECT")
			mqtt_address.Enable()
			mqtt_user.Enable()
			mqtt_password.Enable()

		}

	})
	// create horizontal layout with results
	mqtt_layout := container.New(layout.NewGridLayout(2),
		widget.NewLabel("Address"), mqtt_address,
		widget.NewLabel("Port"), mqtt_port,
		widget.NewLabel("User"), mqtt_user,
		widget.NewLabel("Password"), mqtt_password,
	)

	// return mqtt_layout

	return container.New(layout.NewVBoxLayout(),
		widget.NewLabel("MQTT Broker"),
		brokerSelect,
		mqtt_layout,
		conn_button,
		mqtt_log,
	)
}

func MQTTPubScreen(win fyne.Window) fyne.CanvasObject {

	mqtt_topic := widget.NewEntry()
	mqtt_topic.SetPlaceHolder("devices/<organization>/device/<identifier>/measurements")

	mqtt_message_label := widget.NewLabel("Message:")
	mqtt_message_input := widget.NewMultiLineEntry()
	mqtt_message_input.SetPlaceHolder(`{"temp":24, "hum":50}`)

	mqtt_log := widget.NewLabel("")
	pub_button := widget.NewButton("PUBLISH", func() {
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

	publisherSelect = widget.NewSelectEntry([]string{})
	publisherSelect.OnChanged = func(s string) {
		p := currentBroker.GetPublisher(s)
		mqtt_topic.SetText(p.Topic)
		mqtt_message_input.SetText(p.Message)
	}
	// end publisher implementation

	// Save templates to file

	pLabel := widget.NewLabel("Add name and description")
	namePubEntry := widget.NewEntry()
	namePubEntry.SetPlaceHolder("Name")
	descPubEntry := widget.NewEntry()
	descPubEntry.SetPlaceHolder("Description")
	descPubEntry.MultiLine = true
	cp := container.NewVBox(pLabel, namePubEntry, descPubEntry)

	// Save Button
	addPubButton := widget.NewButton("SAVE TEMPLATE", func() {
		dialog.ShowCustomConfirm("Save Template", "SAVE", "CANCEL", cp, func(b bool) {
			n := model.PublisherData{
				Name:        namePubEntry.Text,
				Description: descPubEntry.Text,
				Topic:       mqtt_topic.Text,
				Message:     mqtt_message_input.Text,
			}
			currentBroker.AddPublisher(n) // add publisher to broker
			config.WriteBroker(currentBroker)
			publisherSelect.SetText(n.Name)
			publisherSelect.SetOptions(currentBroker.GetPublisherNames()) //Add new list to selectEntry
			namePubEntry.SetText("")
			descPubEntry.SetText("")

		}, win)
	})

	// Delete Button
	deletePubButton := widget.NewButton("DELETE TEMPLATE", func() {
		dialog.ShowCustomConfirm("Delete Template", "SAVE", "CANCEL", widget.NewLabel("Are you sure?"), func(b bool) {
			currentBroker.DeletePublisher(publisherSelect.Text)
			config.WriteBroker(currentBroker)
			mqtt_topic.SetText("")         //reset topic entry
			mqtt_message_input.SetText("") //reset message entry
			publisherSelect.SetText("")
			publisherSelect.SetOptions(currentBroker.GetPublisherNames()) //Add new list to selectEntry //todo:change
		}, win)
	})

	pubActions := container.NewHBox(widget.NewLabel("Select Publishing Template:"), addPubButton, deletePubButton)

	return container.New(layout.NewVBoxLayout(),
		widget.NewLabel("Topic:"),
		mqtt_topic,
		mqtt_message_label,
		mqtt_message_input,
		pub_button,
		pubActions,
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
	// containerList.Resize(fyne.NewSize(800, 100))
	index := 1

	scroll := container.NewScroll(containerList)

	desub_button = widget.NewButton("UNSUBSCRIBE", func() {
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

	subs_button = widget.NewButton("SUBSCRIBE", func() {
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
						msgChan <- msg //send msg to channel
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

func testSub(win fyne.Window) fyne.CanvasObject {
	format := false
	chHist := ui.NewChatScrollHistory()
	go func() {
		for msg := range msgChan {
			m := string(msg.Payload())
			if format {
				//pretty json format
				m = string(pretty.Pretty(msg.Payload()))
			}
			console_text := fmt.Sprintf("[%d] topic: %s\ntimestamp:%s\nmsg: %s", chHist.GetLength(), string(msg.Topic()), time.Now().String(), m)
			chHist.AddMessage(console_text)
			fmt.Println(string(msg.Payload()))
		}

	}()

	btn := widget.NewButton("Clear", func() {
		chHist.Clear()
	})

	cn := widget.NewCheck("Pretty JSON", func(b bool) {
		format = b
	})

	a0 := container.New(layout.NewHBoxLayout(), cn, btn)
	a1 := container.New(layout.NewMaxLayout(), chHist.Scroll)
	a2 := container.New(layout.NewBorderLayout(nil, a0, nil, nil), a1, a0)

	return a2
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
	w := a.NewWindow("TGC MQTT Tester v0.5")

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

	pubItem = container.NewTabItemWithIcon("Publish", theme.MailSendIcon(), MQTTPubScreen(w))
	subItem = container.NewTabItemWithIcon("Subscribe", theme.SearchIcon(), MQTTSubScreen(w))
	logItem = container.NewTabItemWithIcon("Log", theme.FileIcon(), testSub(w))

	AppTabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Config", theme.SettingsIcon(), MQTTSetting(w)),
		pubItem,
		subItem,
		logItem,
		//container.NewTabItemWithIcon("Test", theme.ComputerIcon(), makeUI()),
	)

	AppTabs.SetTabLocation(container.TabLocationLeading)
	w.SetContent(AppTabs)
	// w.SetFixedSize(true)
	w.Resize(fyne.NewSize(900, 320))
	w.ShowAndRun()
}
