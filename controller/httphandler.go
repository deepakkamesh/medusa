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
	http.HandleFunc("/api/relayconfigmode", c.relayConfigMode)
	http.HandleFunc("/api/boardconfig", c.boardConfig)
	return http.ListenAndServe(c.httpPort, nil)
}

func (c *Controller) action(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeErr(w, "Err parsing form")
		return
	}

	addr := parseStr(strings.TrimSpace(r.Form.Get("addr")))
	act := strings.TrimSpace(r.Form.Get("actionID"))
	actionID, _ := strconv.ParseUint(act, 16, 8)
	data := parseStr(strings.TrimSpace(r.Form.Get("data")))

	if err := c.core.Action(addr, byte(actionID), data); err != nil {
		writeErr(w, err.Error())
		return
	}
	writeData(w, "ok")
}

func (c *Controller) relayConfigMode(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeErr(w, "Err parsing form")
		return
	}

	hwaddr := parseStr(strings.TrimSpace(r.Form.Get("hwaddr")))
	on, _ := strconv.ParseBool(strings.TrimSpace(r.Form.Get("on")))

	if err := c.core.RelayConfigMode(hwaddr, on); err != nil {
		writeErr(w, err.Error())
		return
	}
	writeData(w, "ok")
}

func (c *Controller) boardConfig(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeErr(w, "Err parsing form")
		return
	}

	hwaddr := parseStr(strings.TrimSpace(r.Form.Get("hwaddr")))
	addr := parseStr(strings.TrimSpace(r.Form.Get("addr")))
	paddr := parseStr(strings.TrimSpace(r.Form.Get("paddr")))
	naddr := parseStr(strings.TrimSpace(r.Form.Get("naddr")))

	if err := c.core.BoardConfig(addr, paddr, hwaddr, naddr); err != nil {
		writeErr(w, err.Error())
		return
	}
	writeData(w, "ok")
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

// parseStr converts comma sep. string hex values to slice.
func parseStr(arg string) []byte {
	if arg == "" {
		return nil
	}
	cmds := strings.Split(arg, ",")
	msg := []byte{}
	for i := 0; i < len(cmds); i++ {
		v, _ := strconv.ParseUint(cmds[i], 16, 8)
		msg = append(msg, byte(v))
	}
	return msg
}

// writeErr is wrapper for writeResponse for Error.
func writeErr(w http.ResponseWriter, s string) {
	writeResponse(w, &response{
		Err: s,
	})
}

// writeData is wrapper for writeResponse for Data.
func writeData(w http.ResponseWriter, s string) {
	writeResponse(w, &response{
		Data: s,
	})
}
