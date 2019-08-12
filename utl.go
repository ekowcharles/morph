package morph

import (
	"fmt"
	"log"
	"os"
	"reflect"
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

func createTypeRegistry(ifc []interface{}) map[string]reflect.Type {
	log.Println("Creating type registry")

	typeRegistry := make(map[string]reflect.Type)

	for _, v := range ifc {
		r := fmt.Sprintf("%T", v)
		log.Printf("Creating registry entry for %s\n", r)

		typeRegistry[r] = reflect.TypeOf(v)
	}

	return typeRegistry
}
