package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func handlePayload(key string, payload []byte) {
	fmt.Print("START PROCESSING MESSAGE")
	switch key {
	case SignInKey:
		fmt.Printf("******* - %v\n\n", payload)
		var signInDto SignInDto
		err := json.Unmarshal(payload, &signInDto)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("******************* SIGN IN *****************")
		fmt.Printf("MESSAGE FROM QUEUE - %s\n", key)
		fmt.Printf("payload - %v\n\n", signInDto)
	case SignUpKey:
		var signUpDto SighUpDto
		err := json.Unmarshal(payload, &signUpDto)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("******************* SIGN UP *****************")
		fmt.Printf("MESSAGE FROM QUEUE - %v", key)
	default:
		log.Panic("handlePayload: unknown payload name")
	}
}
