package main

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/errogaht/bigscreen-tools/repo"
	"github.com/jackc/pgx/v4"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"time"
)

func getTelegramToken() string {
	content, err := os.ReadFile("bot_token.txt")
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

/*func menu(bsRef *bs.Bigscreen) {
	var enteredCommand string
	fmt.Println("Enter command (rooms, participants):")
	fmt.Scan(&enteredCommand)

	switch enteredCommand {
	case "rooms":
		bsRef.PrintOnlineRooms()
			case "participants":
			participants()
	case "exit":
		os.Exit(0)
	}

	menu(bsRef)
}*/

func debug(i interface{}) {

}

func getConn() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
func main() {
	bigscreen := &(bs.Bigscreen{
		JWT: bs.JWTToken{
			Refresh: os.Getenv("BS_JWT_REFRESH"),
			Token:   "renew",
		},
		Bearer:       os.Getenv("BS_BEARER"),
		HostAccounts: os.Getenv("BS_HOST_ACC"),
		HostRealtime: os.Getenv("BS_HOST_REALTIME"),
		Credentials: bs.LoginCredentials{
			Email:    os.Getenv("BS_EMAIL"),
			Password: os.Getenv("BS_PWD"),
		},
		DeviceInfo: fmt.Sprintf(`{"deviceUniqueIdentifier":"%s","drmSystem":"","version":"0.903.19.f05e4d-beta-class-beta","deviceName":"Oculus Quest 2","deviceModel":"Oculus Quest","operatingSystem":"Android OS 10 / API-29 (QQ3A.200805.001/22310100587300000)","CPU":"ARM64 FP ASIMD AES","memory":5842,"GPU":"Adreno (TM) 650"}`, os.Getenv("BS_DEVICE_ID")),
	})

	var rooms []bs.Room
	conn := getConn()
	defer conn.Close(context.Background())

	roomRepo := repo.Room{Conn: conn}
	oculusProfilesRepo := repo.OculusProfile{Conn: conn}
	steamProfilesRepo := repo.SteamProfile{Conn: conn}
	accountProfilesRepo := repo.AccountProfile{Conn: conn}

	for {
		fmt.Println("---------------------------------------------")
		fmt.Printf("%s: start\n", time.Now().Format("2006-01-02 15:04:05"))
		bigscreen.Verify()
		rooms = bigscreen.GetRooms()
		fmt.Printf("%s: got %d rooms\n", time.Now().Format("2006-01-02 15:04:05"), len(rooms))
		//for i := range rooms {
		//	room := &rooms[i]
		//	rooms[i] = bigscreen.GetRoom(room.Id)
		//}

		//debug(rooms)

		oculusProfiles := oculusProfilesRepo.GetProfilesFrom(&rooms)
		oculusProfilesRepo.InsertOrUpdate(&oculusProfiles)
		fmt.Printf("%s: %d oculusProfiles upsert\n", time.Now().Format("2006-01-02 15:04:05"), len(oculusProfiles))

		steamProfiles := steamProfilesRepo.GetProfilesFrom(&rooms)
		steamProfilesRepo.InsertOrUpdate(&steamProfiles)
		fmt.Printf("%s: %d steamProfiles upsert \n", time.Now().Format("2006-01-02 15:04:05"), len(steamProfiles))

		creatorProfiles := accountProfilesRepo.GetCreatorProfilesFrom(&rooms)
		accountProfilesRepo.InsertOrUpdate(&creatorProfiles)

		roomRepo.DeleteAll()
		roomRepo.Insert(&rooms)
		fmt.Printf("%s: rooms refreshed \n", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(30 * time.Second)
	}

	//verify(bigscreen)
	//menu()

	//tgCtx := tg.Context{
	//		Token: "",
	//	}
	//	getMessageTgLoop(&tgCtx, bigscreen)
}

/*func getMessageTgLoop(tgCtxRef *tg.Context, bgCtxRef *bs.Bigscreen) {
	tgCtx := *tgCtxRef

	bot, err := tgbotapi.NewBotAPI(tgCtx.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	updates := bot.GetUpdatesChan(u)
	msgLimit := 4096
	var messages []string
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "rooms" {

			}

			roomsText := bgCtxRef.GetOnlineRooms()

			if len(roomsText) > msgLimit {
				lines := strings.Split(roomsText, "\n")
				var buf string
				for _, line := range lines {
					if len(buf+line) > msgLimit {
						messages = append(messages, buf)
						buf = ""
					}
					buf += line
				}
				messages = append(messages, buf)
			} else {
				messages = append(messages, roomsText)
			}
			for _, message := range messages {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				bot.Send(msg)
			}
		}
	}
}
*/
/*func main() {

	token := de.Login()
	de.Verify(token)
	de.Verify(token)
	de.Verify(token)
	de.Verify(token)
	de.Verify(token)
	de.Verify(token)
	de.LeaveRoom(token)

}*/
