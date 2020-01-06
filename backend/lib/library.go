package lib

import (
  "log"
)

func MyLog(message string) {
	log.Printf("MyLog invoked - printing %s", message)
}
