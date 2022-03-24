package main

import (
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

func getTelegramToken() string {
	content, err := os.ReadFile("bot_token.txt")
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func menu(bsRef *bs.Bigscreen) {
	var enteredCommand string
	fmt.Println("Enter command (rooms, participants):")
	fmt.Scan(&enteredCommand)

	switch enteredCommand {
	case "rooms":
		bsRef.PrintOnlineRooms()
		/*	case "participants":
			participants()*/
	case "exit":
		os.Exit(0)
	}

	menu(bsRef)
}

func debug(i interface{}) {

}

func main() {
	bigscreen := &(bs.Bigscreen{
		JWT:          bs.JWTToken{},
		Bearer:       os.Getenv("BS_BEARER"),
		HostAccounts: os.Getenv("BS_HOST_ACC"),
		HostRealtime: os.Getenv("BS_HOST_REALTIME"),
		Credentials: bs.LoginCredentials{
			Email:    os.Getenv("BS_EMAIL"),
			Password: os.Getenv("BS_PWD"),
		},
		DeviceInfo: fmt.Sprintf(`{"deviceUniqueIdentifier":"%s","drmSystem":"","version":"0.903.19.f05e4d-beta-class-beta","deviceName":"Oculus Quest 2","deviceModel":"Oculus Quest","operatingSystem":"Android OS 10 / API-29 (QQ3A.200805.001/22310100587300000)","CPU":"ARM64 FP ASIMD AES","memory":5842,"GPU":"Adreno (TM) 650"}`, os.Getenv("BS_DEVICE_ID")),
	})

	//rooms := bs.GetRooms(bigscreen)

	//bs.LeaveRoom(bigscreen)
	//participants := bs.Participants(rooms[0].RoomId, bigscreen)
	//debug(participants)
	bigscreen.PrintOnlineRooms()
	//verify(bigscreen)
	//menu()

}
