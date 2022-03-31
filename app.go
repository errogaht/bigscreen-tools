package main

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/errogaht/bigscreen-tools/repo"
	"github.com/errogaht/bigscreen-tools/s"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const SETTING_ROOMS_LAST_UPDATED = "rooms.last.updated"

func getTelegramToken() string {
	content, err := os.ReadFile("bot_token.txt")
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}
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
func logM(m string) {
	fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), m)
}
func roomLoop(bigscreen *bs.Bigscreen) {
	logM("start")
	var rooms []bs.Room
	conn := getConn()
	defer conn.Close(context.Background())
	settingsRepo := repo.Settings{Conn: conn}

	roomRepo := repo.Room{Conn: conn}
	oculusProfilesRepo := repo.OculusProfile{Conn: conn}
	steamProfilesRepo := repo.SteamProfile{Conn: conn}
	accountProfilesRepo := repo.AccountProfile{Conn: conn}

	for {
		fmt.Println("---------------------------------------------")
		bigscreen.Verify()
		rooms = bigscreen.GetRooms()
		logM(fmt.Sprintf("got %d rooms\n", len(rooms)))

		//for i := range rooms {
		//	room := &rooms[i]
		//	rooms[i] = bigscreen.GetRoom(room.Id)
		//}

		//debug(rooms)

		oculusProfiles := oculusProfilesRepo.GetProfilesFrom(&rooms)
		oculusProfilesRepo.Upsert(&oculusProfiles)
		logM(fmt.Sprintf("%d oculusProfiles upsert\n", len(oculusProfiles)))

		steamProfiles := steamProfilesRepo.GetProfilesFrom(&rooms)
		steamProfilesRepo.Upsert(&steamProfiles)
		logM(fmt.Sprintf("%d steamProfiles upsert \n", len(steamProfiles)))

		creatorProfiles := accountProfilesRepo.GetCreatorProfilesFrom(&rooms)
		accountProfilesRepo.Upsert(&creatorProfiles)

		roomRepo.DeleteAll()
		roomRepo.Insert(&rooms)
		settingsRepo.Upsert(&[]s.Settings{{Id: SETTING_ROOMS_LAST_UPDATED, Timestamp: time.Now()}})
		logM(fmt.Sprintf("rooms refreshed, 30s. sleep...\n"))

		time.Sleep(30 * time.Second)
	}
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
		TgToken:      os.Getenv("TG_TOKEN"),
		Credentials: bs.LoginCredentials{
			Email:    os.Getenv("BS_EMAIL"),
			Password: os.Getenv("BS_PWD"),
		},
		DeviceInfo: fmt.Sprintf(`{"deviceUniqueIdentifier":"%s","drmSystem":"","version":"0.903.19.f05e4d-beta-class-beta","deviceName":"Oculus Quest 2","deviceModel":"Oculus Quest","operatingSystem":"Android OS 10 / API-29 (QQ3A.200805.001/22310100587300000)","CPU":"ARM64 FP ASIMD AES","memory":5842,"GPU":"Adreno (TM) 650"}`, os.Getenv("BS_DEVICE_ID")),
	})
	args := os.Args[1:]
	var command string
	if len(args) == 0 {
		command = "help"
	} else {
		command = args[0]
	}
	switch command {
	case "help":
		fmt.Println("Enter command (rooms, participants):")
	case "roomloop":
		roomLoop(bigscreen)
	case "tghook":
		tgHook(bigscreen)
	case "exit":
		os.Exit(0)
	}
}

func tgHook(bgCtxRef *bs.Bigscreen) {
	bot, err := tgbotapi.NewBotAPI(bgCtxRef.TgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook(os.Getenv("TG_WEBHOOK_ROUTE"))
	go http.ListenAndServe("0.0.0.0:8080", nil)

	msgLimit := 4096
	var messages []string

	conn := getConn()
	defer conn.Close(context.Background())
	settingsRepo := repo.Settings{Conn: conn}
	roomsRepo := repo.Room{Conn: conn}
	for update := range updates {
		if update.Message.Text == "rooms" {
			rooms := roomsRepo.FindAll()
			lastUpdated := fmt.Sprintf("Last updated %v ago", time.Now().Sub(settingsRepo.Find(SETTING_ROOMS_LAST_UPDATED).Timestamp))
			roomsText := bgCtxRef.GetOnlineRoomsText(rooms)
			roomsText += "\n" + lastUpdated
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		bot.Send(msg)
		log.Printf("%+v\n", update)
	}

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
