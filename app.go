package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/errogaht/bigscreen-tools/common"
	"github.com/errogaht/bigscreen-tools/repo"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	_ "github.com/joho/godotenv/autoload"
	"html/template"
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

		//i've stuck here because can't achieve stable working GET /rooms/room:id endpoint, i got 500...
		//need to continue hacking bigscreen... everything is ready, just find how to make requests there
		//for i := range rooms {
		//	rooms[i] = bigscreen.GetRoom(rooms[i].Id)
		//}

		//debug(rooms)
		roomRepo.RefreshRoomsInDB(&rooms)

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

	/*conn := NewConn()
	defer conn.Close(context.Background())
	roomsRepo := InitializeRoomRepo()
	//rooms := roomsRepo.FindBy("")
	rooms := bigscreen.GetRooms()
	roomsRepo.RefreshRoomsInDB(&rooms)
	bigscreen.Verify()
	for i := range rooms {
		room := &rooms[i]
		rooms[i] = bigscreen.GetRoom(room.Id)
	}
	debug(rooms)*/

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
	case "tgwebapp":
		tgWebapp(bigscreen)
	case "exit":
		os.Exit(0)
	}
}

func tgWebapp(bgCtxRef *bs.Bigscreen) {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/tgwebapp", func(c *gin.Context) {
		c.HTML(http.StatusOK, "tgwebapp.html", gin.H{
			"title": "Main website",
		})
	})
	router.Run("0.0.0.0:8080")
}

type WebApp struct {
	URL string `json:"url"`
}
type InlineKeyboard struct {
	Text   string `json:"text"`
	WebApp WebApp `json:"web_app"`
}
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboard `json:"inline_keyboard"`
}

func formatAsDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 MST")
}

func getAvatarUrl(p bs.AccountProfile) string {
	if p.OculusProfile.Id != "" {
		if p.OculusProfile.SmallImageURL != "" {
			return p.OculusProfile.SmallImageURL
		}
		if p.OculusProfile.ImageURL != "" {
			return p.OculusProfile.ImageURL
		}
	}

	if p.SteamProfile.Id != "" {
		if p.SteamProfile.Avatar != "" {
			return p.SteamProfile.Avatar
		}
	}

	return ""
}

func tgHook(bgCtxRef *bs.Bigscreen) {
	bot, err := tgbotapi.NewBotAPI(bgCtxRef.TgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

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

	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
		"getAvatarUrl": getAvatarUrl,
	})
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")

	router.GET("/tgwebapp", func(c *gin.Context) {
		rooms := roomsRepo.FindBy("")
		marshal, err := json.Marshal(rooms)
		if err != nil {
			return
		}
		c.HTML(http.StatusOK, "bootstrap.gohtml", gin.H{
			"Rooms":     rooms,
			"roomsJson": template.JS(marshal),
		})
	})
	router.POST(os.Getenv("TG_WEBHOOK_ROUTE"), func(c *gin.Context) {
		update, err := bot.HandleUpdate(c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else if update.Message == nil {
			c.JSON(200, gin.H{})
		} else if !update.Message.IsCommand() {
			c.JSON(200, gin.H{})
		} else {
			switch update.Message.Command() {

			case "help":
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I understand /menu, /all, /chat, /movies, /gaming, /nsfw, /sports")
				msg.ReplyMarkup = menuKeyboard
				msg.DisableNotification = true
				bot.Send(msg)
			case "all":
				rooms := roomsRepo.FindBy("")
				sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, update, bot)
			case "chat":
				rooms := roomsRepo.FindBy("category = $1", "CHAT")
				sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, update, bot)
			case "movies":
				rooms := roomsRepo.FindBy("category = $1", "MOVIES")
				sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, update, bot)
			case "gaming":
				rooms := roomsRepo.FindBy("category = $1", "GAMING")
				sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, update, bot)
			case "sports":
				rooms := roomsRepo.FindBy("category = $1", "SPORTS")
				sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, update, bot)
			case "nsfw":
				rooms := roomsRepo.FindBy("category = $1", "NSFW")
				sendTgRoomsMessages(rooms, bgCtxRef, settingsRepo, update, bot)
			case "menu":
				chat := update.FromChat()
				if !chat.IsPrivate() {
					msg := tgbotapi.NewMessage(
						update.Message.Chat.ID,
						fmt.Sprintf("This cool, convenient menu works only in private chats. Open chat with me @%s and enjoy the feature!", bot.Self.UserName),
					)
					msg.DisableNotification = true
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Press button bellow")
					msg.DisableNotification = true
					msg.ReplyMarkup = InlineKeyboardMarkup{
						InlineKeyboard: [][]InlineKeyboard{
							{InlineKeyboard{
								Text:   "Open menu",
								WebApp: WebApp{URL: "https://xxx/yyy"},
							}},
						},
					}
					bot.Send(msg)
				}
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
				msg.DisableNotification = true
				bot.Send(msg)
			}

			c.JSON(200, gin.H{})
		}
	})
	router.Run("0.0.0.0:8080")
}
func sendTgRoomsMessages(rooms *[]bs.Room, bgCtxRef *bs.Bigscreen, settingsRepo *repo.SettingsRepo, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msgLimit := 4096
	var messages []string
	lastUpdated := fmt.Sprintf("Last updated %.0fs ago", time.Now().Sub(settingsRepo.Find(repo.SETTING_ROOMS_LAST_UPDATED).Timestamp).Seconds())
	if len(*rooms) == 0 {
		msgText := "No rooms found.\n"
		msgText += lastUpdated
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		msg.DisableNotification = true
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
		msg.DisableNotification = true
		bot.Send(msg)
	}
}
