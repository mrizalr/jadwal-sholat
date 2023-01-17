package city

import (
	"encoding/json"
	"fmt"

	"github.com/mrizalr/jadwal-sholat/network"
	"github.com/spf13/viper"
)

type GetCityResponse struct {
	Status bool `json:"status"`
	Data   []struct {
		Id     string `json:"id"`
		Lokasi string `json:"lokasi"`
	} `json:"data"`
}

func fetchCities() (GetCityResponse, error) {
	result := GetCityResponse{}

	url := fmt.Sprintf("https://api.myquran.com/v1/sholat/kota/cari/%s", viper.GetString("kota"))
	response, err := network.FetchData("GET", url, nil)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(response, &result)
	return result, err
}

func setCityID(cities GetCityResponse) {
	var cityID string
	if len(cities.Data) > 1 {
		citiesOption := ""
		for idx, value := range cities.Data {
			idKota := viper.GetString("id_kota")
			if idKota != "" && idKota == value.Id {
				return
			}

			citiesOption += fmt.Sprintf("%d. %s - %s\n", idx+1, value.Id, value.Lokasi)
		}
		fmt.Println(citiesOption)
		fmt.Printf("Masukkan kode kota yang anda pilih : ")
		fmt.Scanln(&cityID)
	} else {
		cityID = cities.Data[0].Id
	}

	viper.Set("id_kota", cityID)
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}
}

func SetCity() {
	cities, err := fetchCities()
	if err != nil {
		panic(err)
	}

	if cities.Status == false {
		panic(fmt.Errorf("error when fetch cities data : %s", err.Error()))
	}

	setCityID(cities)
}
