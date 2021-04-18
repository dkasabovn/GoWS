package config

import (
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FsClient *firestore.Client

func CreateFirestoreClient() {
	sa := option.WithCredentialsFile("./serviceAccount.json")
	conf := &firebase.Config{ProjectID: "compete-fd1d3"}
	app, err := firebase.NewApp(CTX, conf, sa)
	if err != nil {
		panic(err)
	}

	client, err := app.Firestore(CTX)
	if err != nil {
		panic(err)
	}
	FsClient = client
}

func Close() {
	FsClient.Close()
}
