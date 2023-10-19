package api

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go/v5"
	"oursos.com/packages/util"
)

func CreateBroadcast() {
	enverr := godotenv.Load()
	util.CheckError(enverr)

	pusherClient := pusher.Client{
		AppID:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: "us3",
		Secure:  true,
	}

	data := map[string]string{"message": "Special Alert"}

	err := pusherClient.Trigger("my-channel", "my-event", data)
	util.CheckError(err)
}
