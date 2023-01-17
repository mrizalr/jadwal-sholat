package schedule

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/mrizalr/jadwal-sholat/network"
	"github.com/spf13/viper"
)

type ScheduleResponse struct {
	Status bool `json:"status"`
	Data   struct {
		ID        string `json:"id"`
		Lokasi    string `json:"lokasi"`
		Daerah    string `json:"daerah"`
		Koordinat struct {
			Lat     float64 `json:"lat"`
			Lon     float64 `json:"lon"`
			Lintang string  `json:"lintang"`
			Bujur   string  `json:"bujur"`
		} `json:"koordinat"`
		Jadwal struct {
			Tanggal string `json:"tanggal"`
			Imsak   string `json:"imsak"`
			Subuh   string `json:"subuh"`
			Terbit  string `json:"terbit"`
			Dhuha   string `json:"dhuha"`
			Dzuhur  string `json:"dzuhur"`
			Ashar   string `json:"ashar"`
			Maghrib string `json:"maghrib"`
			Isya    string `json:"isya"`
			Date    string `json:"date"`
		} `json:"jadwal"`
	} `json:"data"`
}

type JadwalSholat struct {
	Imsak   time.Time `json:"imsak"`
	Subuh   time.Time `json:"subuh"`
	Terbit  time.Time `json:"terbit"`
	Dhuha   time.Time `json:"dhuha"`
	Dzuhur  time.Time `json:"dzuhur"`
	Ashar   time.Time `json:"ashar"`
	Maghrib time.Time `json:"maghrib"`
	Isya    time.Time `json:"isya"`
}

func fetchSchedule() (ScheduleResponse, error) {
	cityID := viper.GetString("id_kota")
	date := time.Now()

	var result ScheduleResponse
	url := fmt.Sprintf("https://api.myquran.com/v1/sholat/jadwal/%s/%d/%d/%d", cityID, date.Year(), date.Month(), date.Day())
	response, err := network.FetchData("GET", url, nil)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(response, &result)
	return result, err
}

func FindTodayJadwal() JadwalSholat {
	schedule := getSchedule()

	result := JadwalSholat{}
	resultFields := reflect.ValueOf(&result)

	jadwalData := schedule.Data.Jadwal
	dataFields := reflect.ValueOf(jadwalData)

	for i := 1; i < dataFields.NumField()-1; i++ {
		val := dataFields.Field(i).String()
		clock, _ := time.Parse("15:04", val)

		date := time.Now()
		jadwal := time.Date(date.Year(), date.Month(), date.Day(), clock.Hour(), clock.Minute(), 0, 0, time.Local)
		resultFields.Elem().FieldByName(dataFields.Type().Field(i).Name).Set(reflect.ValueOf(jadwal))
	}
	return result
}

func getSchedule() ScheduleResponse {
	schedule, err := fetchSchedule()
	if err != nil {
		panic(fmt.Errorf("error when fetch schedule data : %s", err.Error()))
	}

	return schedule
}

func FindNextJadwal(jadwalSholat JadwalSholat) int {
	jadwalFields := reflect.ValueOf(jadwalSholat)
	for i := 0; i < jadwalFields.NumField(); i++ {
		sholatTime := jadwalFields.Field(i).Interface().(time.Time)
		if time.Now().After(sholatTime) == false && viper.GetBool(fmt.Sprintf("notify.%s", jadwalFields.Type().Field(i).Name)) == true {
			return i
		}
	}
	return -1
}
