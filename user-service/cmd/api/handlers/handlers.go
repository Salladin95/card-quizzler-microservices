package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/cachedRepository"
	user "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/entities"
	userService "github.com/Salladin95/card-quizzler-microservices/user-service/proto"
	"log"
)

type UserServer struct {
	userService.UnimplementedUserServiceServer
	Repo cachedRepository.CachedRepository
}

func HandleRabbitPayload(key string, payload []byte) {
	fmt.Print("START PROCESSING MESSAGE")
	switch key {
	case constants.SignInKey:
		fmt.Printf("******* - %v\n\n", payload)
		var signInDto user.SignInDto
		err := json.Unmarshal(payload, &signInDto)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("******************* SIGN IN *****************")
		fmt.Printf("MESSAGE FROM QUEUE - %s\n", key)
		fmt.Printf("payload - %v\n\n", signInDto)
	case constants.SignUpKey:
		var signUpDto user.SignUpDto
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
