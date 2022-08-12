package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// response Struct to return JSON.
type response struct {
	Err  string
	Data interface{}
}

func main() {

	params := url.Values{
		"cmd": {os.Args[3]},
	}

	resp, err := http.PostForm("http://"+os.Args[2]+"/api/cli", params)
	if err != nil {
		fmt.Printf("request failed: %v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("http read failed: %v", err)
		return
	}
	respo := &response{}
	if err := json.Unmarshal(body, respo); err != nil {
		fmt.Printf("json unmarshal failed: %v", err)
		return
	}
	if respo.Err != "" {
		fmt.Printf("Core error: %s", respo.Err)
		return
	}
	fmt.Println(respo.Data)
}
