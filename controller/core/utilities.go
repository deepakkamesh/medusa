package core

import "fmt"

// PrintPkt returns a string of packet formatted properly.
func PrintPkt(preamble string, b []byte, s int) string {
	logMsg := preamble
	for i := 0; i < s; i++ {
		logMsg = logMsg + fmt.Sprintf("%X ", b[i])
	}
	return logMsg
}
