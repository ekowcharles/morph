package morph

import (
	"fmt"
	"log"
	"os"
)

const zeroString = ""

func panicIf(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func createTypeRegistry(ifc []interface{}) map[string]interface{} {
	log.Println("Creating type registry")

	typeRegistry := make(map[string]interface{})

	for _, v := range ifc {
		r := fmt.Sprintf("%T", v)
		log.Printf("Creating registry entry for %s\n", r)

		typeRegistry[r] = v
	}

	return typeRegistry
}
