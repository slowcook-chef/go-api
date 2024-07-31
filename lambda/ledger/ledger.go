package ledger

import (
	"log"
)

func LogError(err *error) {
	log.Println("ERROR:", err)
}

func LogHandlerStart(message string) {
	log.Printf("HANDLER START log: %s", message)
}

func LogHandlerEnd(message string) {
	log.Printf("HANDLER END log: %s", message)
}

func LogHandlerProcess(message string) {
	log.Printf("PROCESS log: %s", message)
}
