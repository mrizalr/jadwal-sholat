package main

import (
	"fmt"
	"reflect"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/mrizalr/jadwal-sholat/city"
	"github.com/mrizalr/jadwal-sholat/config"
	"github.com/mrizalr/jadwal-sholat/schedule"
	"github.com/spf13/viper"
)

var todayJadwal schedule.JadwalSholat
var a fyne.App

type sholatInfo struct {
	nextTime time.Time
	nextName string
}

func init() {
	config.InitConfig()
	city.SetCity()
	todayJadwal = schedule.FindTodayJadwal()
}

func main() {
	a = app.New()

	jadwalFields := reflect.ValueOf(todayJadwal)
	jadwalInfo := ""
	for i := 0; i < jadwalFields.NumField(); i++ {
		sholatName := jadwalFields.Type().Field(i).Name
		sholatTime := jadwalFields.Field(i).Interface().(time.Time).Format("15:04")
		jadwalInfo += fmt.Sprintf("%s - %s\n", sholatName, sholatTime)
	}

	w := CreateWindow(a, "Jadwal sholat", jadwalInfo)
	w.Show()

	go func() {
		info := getNextSholatInfo(jadwalFields)
		fmt.Println(info.nextTime)

		alertDone := map[int]bool{0: false, 5: false, 10: false, 20: false}

		for tick := range time.Tick(time.Second) {
			diff := info.nextTime.Sub(tick)

			if tick.After(info.nextTime) && viper.GetBool("period.0") && alertDone[0] == false {
				w := CreateWindow(a, "Saatnya sholat", info.nextName)
				w.Show()

				info = getNextSholatInfo(jadwalFields)
				for k, _ := range alertDone {
					alertDone[k] = false
				}
			} else if diff < time.Duration(5*time.Minute) && viper.GetBool("period.5") && alertDone[5] == false {
				w := CreateWindow(a, "5 Menit lagi", info.nextName)
				w.Show()
				alertDone[5] = true
			} else if diff < time.Duration(10*time.Minute) && viper.GetBool("period.10") && alertDone[10] == false {
				w := CreateWindow(a, "10 Menit lagi", info.nextName)
				w.Show()
				alertDone[10] = true
			} else if diff < time.Duration(20*time.Minute) && viper.GetBool("period.20") && alertDone[20] == false {
				w := CreateWindow(a, "20 Menit lagi", info.nextName)
				w.Show()
				alertDone[20] = true
			}
		}
	}()
	a.Run()
}

func getNextSholatInfo(jadwalFields reflect.Value) sholatInfo {
	next := schedule.FindNextJadwal(todayJadwal)
	nextSholatTime := jadwalFields.Field(next).Interface().(time.Time)
	nextSholatName := jadwalFields.Type().Field(next).Name
	return sholatInfo{
		nextTime: nextSholatTime,
		nextName: nextSholatName,
	}
}

func CreateWindow(app fyne.App, countdownStr, sholatTypeStr string) fyne.Window {
	w := app.NewWindow("Notify")
	w.Resize(fyne.NewSize(300, 100))

	container := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	countdownLabel := widget.NewLabel(countdownStr)
	sholatType := widget.NewLabel(sholatTypeStr)
	container.Add(countdownLabel)
	container.Add(sholatType)

	w.SetContent(container)
	w.SetFixedSize(true)
	w.CenterOnScreen()

	return w
}
