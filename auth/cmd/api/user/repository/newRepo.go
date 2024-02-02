package user

import (
	"cloud.google.com/go/firestore"
	fireBaseAuth "firebase.google.com/go/v4/auth"
)

func NewUserRepository(dbClient *firestore.Client, authClient *fireBaseAuth.Client) Repository {
	return &repository{dbClient: dbClient, authClient: authClient}
}
