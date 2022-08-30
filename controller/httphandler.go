package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// response Struct to return JSON.
type response struct {
	Err  string      // Error message.
	Data interface{} // Data message.
}

// StartHTTP starts the HTTP server.
func (c *Controller) StartHTTP() error {
	http.HandleFunc("/api/cli", c.cli)
	return http.ListenAndServe(c.httpPort, nil)
}

// cli handles the command line raw packets.
func (c *Controller) cli(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		writeResponse(w, &response{
			Err: string("Error parsing form"),
		})
		return
	}

	// convert string into byte array.
	proto := strings.TrimSpace(r.Form.Get("proto"))
	cmd := strings.TrimSpace(r.Form.Get("cmd"))
	cmds := strings.Split(cmd, ",")
	msg := []byte{}
	strMsg := ""
	for i := 0; i < len(cmds); i++ {
		v, _ := strconv.ParseUint(cmds[i], 16, 8)
		msg = append(msg, byte(v))
		strMsg = strMsg + fmt.Sprintf("%X ", byte(v))
	}

	glog.Infof("Cli Command to %v: %v", proto, strMsg)

	switch {
	// Send to UDP.
	case proto == "U":
		if err := c.core.SendManualRelayCfg(msg); err != nil {
			glog.Errorf("Error sending relay config %v cmd:  %v", strMsg, err)
		}
		// Send to TCP.
	case proto == "T":
		if err := c.core.SendRawPacket(msg); err != nil {
			glog.Errorf("Error sending pkt  %v cmd:  %v", strMsg, err)
		}
	}

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
	w.Write(js)
}
