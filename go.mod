module github.com/dhinojosac/mqtt_app_desktop

go 1.15

replace github.com/dhinojosac/mqtt_app_desktop/logger => ./logger

replace github.com/dhinojosac/mqtt_app_desktop/config => ./config

replace github.com/dhinojosac/mqtt_app_desktop/model => ./model

require (
	fyne.io/fyne v1.4.3
	fyne.io/fyne/v2 v2.0.4
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/spf13/viper v1.8.1
	go.uber.org/zap v1.18.1
)
