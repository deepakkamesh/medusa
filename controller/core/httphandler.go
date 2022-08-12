package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// response Struct to return JSON.
type response struct {
	Err  string
	Data interface{}
}

func (c *Core) StartHTTP() error {

	http.HandleFunc("/api/cli", c.cli)
	return http.ListenAndServe(c.hostPort, nil)
}

// cli handles the command line raw packets.
func (c *Core) cli(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		writeResponse(w, &response{
			Err: string("Error parsing form"),
		})
		return
	}

	// convert string into byte array.
	cmd := strings.TrimSpace(r.Form.Get("cmd"))
	pkt := strings.Split(cmd, ",")
	message := []byte{}
	for i := 0; i < len(pkt); i++ {
		v, _ := strconv.ParseUint(pkt[i], 16, 8)
		message = append(message, byte(v))
	}
	fmt.Println(message)

	// If Exec is successful, send back command output.
	writeResponse(w, &response{
		Data: string("ok"),
	})

}

// writeResponse writes the response json object to w. If unable to marshal
// it writes a http 500.
func writeResponse(w http.ResponseWriter, resp *response) {
	w.Header().Set("Content-Type", "application/json")
	js, e := json.Marshal(resp)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	//	log.Printf("Writing json response %s", js)
	w.Write(js)
}
