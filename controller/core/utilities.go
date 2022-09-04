package core

import "fmt"

// PrintPkt returns a string of packet formatted properly.
func PrintPkt(preamble string, b []byte) string {
	logMsg := preamble
	for i := 0; i < len(b); i++ {
		logMsg = logMsg + fmt.Sprintf("%X ", b[i])
	}
	return logMsg
}

// PrintPkt returns a string of packet formatted properly.
func PP(b []byte, format string, a ...any) string {
	logMsg := fmt.Sprintf(format, a...)
	for i := 0; i < len(b); i++ {
		logMsg = logMsg + fmt.Sprintf("%X ", b[i])
	}
	return logMsg
}

// PrintPkt returns a string of packet formatted properly.
func PP2(b []byte) string {
	logMsg := ""
	for i := 0; i < len(b); i++ {
		logMsg = logMsg + fmt.Sprintf("%X ", b[i])
	}
	return logMsg
}
