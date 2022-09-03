package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true

	app.Commands = []*cli.Command{
		{
			Name:    "led",
			Aliases: []string{"l"},
			Usage:   "LED control",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "a",
					Usage: "Board Address hex -a AB,FF,3A",
				},
				&cli.BoolFlag{
					Name:  "on",
					Usage: "on",
				},
			},
			Action: func(c *cli.Context) error {
				addr := c.String("a")
				on := c.Bool("on")

				onV := string("1")
				if !on {
					onV = "0"
				}

				params := url.Values{
					"addr":     {addr},
					"actionID": {strconv.Itoa(core.ActionLED)},
					"data":     {onV},
				}
				post("", params, "action")

				return nil
			},
		},
		{
			Name:    "temp",
			Aliases: []string{"t"},
			Usage:   "get Temp/Humidity",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "a",
					Usage: "Board Address hex -a AB,FF,3A",
				},
			},
			Action: func(c *cli.Context) error {
				addr := parseStr(c.String("a"))
				fmt.Println(addr)
				return nil
			},
		},
		{
			Name:    "RelayConfigMode",
			Aliases: []string{"rcm"},
			Usage:   "Enable/Disable config mode for relay",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "a",
					Usage: "Relay HWAddress hex -a AB,FF,3A",
				},
				&cli.BoolFlag{
					Name:  "on",
					Usage: "on",
				},
			},
			Action: func(c *cli.Context) error {
				addr := parseStr(c.String("a"))
				on := c.Bool("on")
				fmt.Println(addr, on)
				return nil
			},
		},
		{
			Name:    "BoardConfig",
			Aliases: []string{"bc"},
			Usage:   "Send board config",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "a",
					Usage: "Board Address hex -a AB,FF,3A",
				},
				&cli.StringFlag{
					Name:  "pa",
					Usage: "Pipe Address hex -a AB,FF,3A",
				},
				&cli.StringFlag{
					Name:  "hw",
					Usage: "HWAddr Address hex -a AB,FF,3A",
				},

				&cli.StringFlag{
					Name:  "na",
					Usage: "Board address from config (hex) to send -na A1,22,3",
				},
			},
			Action: func(c *cli.Context) error {
				addr := parseStr(c.String("a"))
				paddr := parseStr(c.String("pa"))
				hwaddr := parseStr(c.String("hw"))
				na := parseStr(c.String("na"))

				fmt.Println(addr, paddr, hwaddr, na)
				return nil
			},
		},
	}
	app.Run(os.Args)

}

// parseStr converts comma sep. string hex values to slice.
func parseStr(arg string) []byte {
	cmds := strings.Split(arg, ",")
	msg := []byte{}
	for i := 0; i < len(cmds); i++ {
		v, _ := strconv.ParseUint(cmds[i], 16, 8)
		msg = append(msg, byte(v))
	}
	return msg
}

// response Struct to return JSON.
type response struct {
	Err  string
	Data interface{}
}

func post(host string, params url.Values, api string) {
	/*	params = url.Values{
		"proto": {os.Args[3]},
		"cmd":   {os.Args[4]},
	}*/

	resp, err := http.PostForm("http://"+host+"/api/"+api, params)
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
	resp.Body.Close()
}
