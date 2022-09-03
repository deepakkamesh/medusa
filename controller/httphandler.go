package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// response Struct to return JSON.
type response struct {
	Err  string      // Error message.
	Data interface{} // Data message.
}

// StartHTTP starts the HTTP server.
func (c *Controller) StartHTTP() error {
	http.HandleFunc("/api/action", c.action)
	return http.ListenAndServe(c.httpPort, nil)
}

// led controls the led.
func (c *Controller) action(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeErr(w, "Err parsing form")
		return
	}

	addr := parseStr(strings.TrimSpace(r.Form.Get("addr")))
	act := strings.TrimSpace(r.Form.Get("actionID"))
	actionID, _ := strconv.ParseUint(act, 16, 8)
	data := parseStr(strings.TrimSpace(r.Form.Get("data")))

	_, _, _ = addr, actionID, data

	//TODO:
	/*
		if err := c.core.LEDOn(addr, on); err != nil {
			writeErr(w, err.Error())
			return
		}*/

	writeData(w, "ok")
}

/*
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

}*/

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

func writeErr(w http.ResponseWriter, s string) {
	writeResponse(w, &response{
		Err: s,
	})
}

func writeData(w http.ResponseWriter, s string) {
	writeResponse(w, &response{
		Data: s,
	})
}
