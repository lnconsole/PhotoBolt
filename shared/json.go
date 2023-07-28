package shared

import (
	"encoding/json"
	"log"
)

func CleanMarshal(obj any) []byte {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Print(err)
	}

	return bytes
}

func CleanUnmarshal(data []byte, obj any) {
	err := json.Unmarshal(data, &obj)
	if err != nil {
		log.Print(err)
	}
}
