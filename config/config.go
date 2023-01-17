package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/viper"
)

func checkConfig() {
	_, err := os.Stat("config.json")
	if err != nil {
		createConfig()
		return
	}
}

func createConfig() {
	file, err := os.Create("config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	configData := struct {
		Kota   string `json:"kota"`
		IDKota string `json:"id_kota"`
		Period map[string]bool
		Notify map[string]bool
	}{
		Kota:   "surabaya",
		IDKota: "",
		Period: map[string]bool{
			"20": true,
			"10": false,
			"5":  true,
			"0":  true,
		},
		Notify: map[string]bool{
			"imsak":   false,
			"subuh":   true,
			"terbit":  false,
			"dhuha":   false,
			"dzuhur":  true,
			"ashar":   true,
			"maghrib": true,
			"isya":    true,
		},
	}

	json, err := json.Marshal(configData)
	if err != nil {
		panic(err)
	}
	io.WriteString(file, string(json))
}

func InitConfig() {
	checkConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	printConfigInfo()
}

func printConfigInfo() {
	fmt.Printf(`
		Config info

		=================================
		Kota : %s
		=================================
		
		Notify

		=================================
		20 : %s
		10 : %s
		5 : %s
		0 : %s
		=================================
		imsak : %s
		subuh : %s
		terbit : %s
		dhuha : %s
		dzuhur : %s
		ashar : %s
		maghrib : %s
		isya : %s
		=================================
		`, viper.GetString("kota"),
		viper.GetString("period.20"),
		viper.GetString("period.10"),
		viper.GetString("period.5"),
		viper.GetString("period.0"),
		viper.GetString("notify.imsak"),
		viper.GetString("notify.subuh"),
		viper.GetString("notify.terbit"),
		viper.GetString("notify.dhuha"),
		viper.GetString("notify.dzuhur"),
		viper.GetString("notify.ashar"),
		viper.GetString("notify.maghrib"),
		viper.GetString("notify.isya"),
	)
}
