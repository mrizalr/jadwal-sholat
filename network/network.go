package network

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func FetchData(method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	try := 3
	client := &http.Client{}
	var resp *http.Response

	for {
		resp, err = client.Do(req)
		try--
		fmt.Println("fetching data ...")

		if err != nil {
			if try > 0 {
				time.Sleep(time.Second)
				fmt.Println("fail fetching data, check your connection")
				continue
			}
			return nil, fmt.Errorf("Cannot fetch data : %s", err.Error())
		} else {
			break
		}
	}

	json, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return json, nil
}
