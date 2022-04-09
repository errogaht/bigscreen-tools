package main

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/errogaht/bigscreen-tools/common"
	"github.com/errogaht/bigscreen-tools/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func getTelegramToken() string {
	content, err := os.ReadFile("bot_token.txt")
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}
func debug(i interface{}) {

}

var conn *pgx.Conn

func NewConn() *pgx.Conn {
	if nil == conn {
		var err error
		conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
	}
	return conn
}

func roomLoop(bigscreen *bs.Bigscreen) {
	common.LogM("start")
	var rooms []bs.Room
	conn := NewConn()
	defer conn.Close(context.Background())

	roomRepo := InitializeRoomRepo()

	for {
		fmt.Println("---------------------------------------------")
		bigscreen.Verify()
		rooms = bigscreen.GetRooms()
		common.LogM(fmt.Sprintf("got %d rooms", len(rooms)))

		//for i := range rooms {
		//	room := &rooms[i]
		//	rooms[i] = bigscreen.GetRoom(room.Id)
		//}

		//debug(rooms)
		roomRepo.RefreshRoomsInDB(&rooms)

		time.Sleep(30 * time.Second)
	}
}
func main() {
	/*conn := NewConn()
	defer conn.Close(context.Background())
	roomsRepo := InitializeRoomRepo()
	rooms := roomsRepo.FindBy("category = $1", "GAMING")
	rooms = roomsRepo.FindBy("category = $1", "CHAT")

	debug(rooms)*/
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

	conn := NewConn()
	defer conn.Close(context.Background())
	settingsRepo := InitializeSettingsRepo()
	roomsRepo := InitializeRoomRepo()

	menuKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/all"),
			tgbotapi.NewKeyboardButton("/chat"),
			tgbotapi.NewKeyboardButton("/movies"),
			tgbotapi.NewKeyboardButton("/gaming"),
			tgbotapi.NewKeyboardButton("/nsfw"),
			tgbotapi.NewKeyboardButton("/sports"),
		),
	)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {

		case "help":
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I understand /all, /chat, /movies, /gaming, /nsfw")
			msg.ReplyMarkup = menuKeyboard
			bot.Send(msg)
		case "all":
			rooms := roomsRepo.FindBy("")
			sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, &update, bot)
		case "chat":
			rooms := roomsRepo.FindBy("category = $1", "CHAT")
			sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, &update, bot)
		case "movies":
			rooms := roomsRepo.FindBy("category = $1", "MOVIES")
			sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, &update, bot)
		case "gaming":
			rooms := roomsRepo.FindBy("category = $1", "GAMING")
			sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, &update, bot)
		case "sports":
			rooms := roomsRepo.FindBy("category = $1", "SPORTS")
			sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, &update, bot)
		case "nsfw":
			rooms := roomsRepo.FindBy("category = $1", "NSFW")
			sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, &update, bot)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
			bot.Send(msg)
		}
	}

}
func sendTgRoomsMessages(rooms *[]bs.Room, bgCtxRef *bs.Bigscreen, settingsRepo *repo.SettingsRepo, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msgLimit := 4096
	var messages []string
	lastUpdated := fmt.Sprintf("Last updated %v ago", time.Now().Sub(settingsRepo.Find(repo.SETTING_ROOMS_LAST_UPDATED).Timestamp))
	if len(*rooms) == 0 {
		msgText := "No rooms found.\n"
		msgText += lastUpdated
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		bot.Send(msg)
		return
	}
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
	for i := range messages {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages[i])
		bot.Send(msg)
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
