package server

import "log"

func logf(format string, v ...any) {
	log.Printf(format, v...)
}
