package bs

import (
	"encoding/json"
	"fmt"
	"gopkg.in/guregu/null.v4"
	"log"
	"os"
	"strings"
	"time"
)

type Timestamp struct {
	Seconds     int64 `json:"_seconds"`
	Nanoseconds int64 `json:"_nanoseconds"`
}

type RoomCreator struct {
	UserSessionId string
	IsMod         bool
	IsStaff       bool
}

type AccountProfile struct {
	CreatedAtTimestamp Timestamp   `json:"createdAt"`
	CreatedAt          time.Time   `db:"created_at"`
	IsVerified         bool        `db:"is_verified"`
	IsBanned           bool        `db:"is_banned"`
	IsStaff            bool        `db:"is_staff"`
	Username           string      `db:"username"`
	SteamProfileId     null.String `db:"steam_profile_id"`
	OculusProfileId    null.String `db:"oculus_profile_id"`
	SteamProfile       SteamProfile
	OculusProfile      OculusProfile
}
type OculusProfile struct {
	Id            string `json:"oculusId"`
	ImageURL      string `json:"oculusImageURL"`
	SmallImageURL string `json:"oculusSmallImageURL"`
}
type SteamProfile struct {
	Id                       string `json:"steamid"`
	CommunityVisibilityState uint8  `json:"communityvisibilitystate"`
	ProfileState             uint8  `json:"profilestate"`
	PersonaName              string `json:"personaname"`
	ProfileUrl               string `json:"profileurl"`
	Avatar                   string `json:"avatar"`
	AvatarMedium             string `json:"avatarmedium"`
	AvatarFull               string `json:"avatarfull"`
	AvatarHash               string `json:"avatarhash"`
	PersonaState             uint8  `json:"personastate"`
	RealName                 string `json:"realname"`
	PrimaryClanId            string `json:"primaryclanid"`
	TimeCreated              uint64 `json:"timecreated"`
	CreatedAt                time.Time
	PersonaStateFlags        uint16 `json:"personastateflags"`
	LocCountryCode           string `json:"loccountrycode"`
}
type RoomUser struct {
	IsAdmin          bool   `db:"is_admin"`
	IsMod            bool   `db:"is_mod"`
	IsStaff          bool   `db:"is_staff"`
	Version          string `db:"version"`
	UserSessionId    string `db:"user_session_id"`
	LegacyUserId     string
	RoomId           string    `db:"room_id"`
	SeatIndex        uint8     `db:"seat_index"`
	CreatedAt        time.Time `db:"created_at"`
	AccountProfileId string    `db:"account_profile"`
	AccountProfile   AccountProfile
}
type Room struct {
	Creator                RoomCreator
	CreatorProfile         AccountProfile
	CreatorProfileUsername string
	Name                   string
	Description            string
	Category               string
	Environment            string
	Size                   uint8
	Version                string
	RoomType               string
	Visibility             string
	InviteCode             string
	Status                 string
	RemoteUsers            []RoomUser
	Participants           uint8
	CreatedAt              time.Time
	RoomId                 string
	Id                     string
}

func (bsRef *Bigscreen) GetRooms() (rooms []Room) {
	body, _ := bsRef.request(
		(*bsRef).HostRealtime+"/rooms/latest",
		"GET",
		make(map[string]string),
		"",
	)

	err := json.Unmarshal(body, &rooms)
	if err != nil {
		log.Panic(err.Error())
	}
	for i := range rooms {
		r := &rooms[i]
		r.Id = strings.Split(r.RoomId, ":")[1]
		r.CreatorProfile.CreatedAt = time.Unix(r.CreatorProfile.CreatedAtTimestamp.Seconds, r.CreatorProfile.CreatedAtTimestamp.Nanoseconds)
	}
	return
}

func getMsgTemplate() string {
	content, err := os.ReadFile("roomsMsgTemplate.txt")
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func (bsRef *Bigscreen) GetOnlineRoomsText(listOfRoomsRef *[]Room) string {
	listOfRooms := *listOfRoomsRef

	var result string
	//9.   CHAT -  5/15 - Steve's Place 21               | U:Steve_               Music
	for i, room := range listOfRooms {
		i++
		result += fmt.Sprintf("%d. %s", i, room.Name)
		if room.Description != "" {
			result += " (" + room.Description + ")"
		}
		result += "\n"

		result += fmt.Sprintf(
			"%s - %d/%d %s\n\n",
			room.Category,
			room.Participants,
			room.Size,
			room.CreatorProfile.Username,
			//room.Id,
		)
	}

	return result
}

func (bsRef *Bigscreen) GetRoom(roomId string) (room Room) {
	body, _ := bsRef.request(
		(*bsRef).HostRealtime+"/room/room:"+roomId,
		"GET",
		make(map[string]string),
		"",
	)

	err := json.Unmarshal(body, &room)
	if err != nil {
		log.Panic(err.Error())
	}
	for i := range room.RemoteUsers {
		u := &room.RemoteUsers[i]
		u.AccountProfile.CreatedAt = time.Unix(u.AccountProfile.CreatedAtTimestamp.Seconds, u.AccountProfile.CreatedAtTimestamp.Nanoseconds)
	}
	return
}

func (bsRef *Bigscreen) LeaveRoom() {
	bsRef.request(
		(*bsRef).HostRealtime+"/leave_room",
		"GET",
		make(map[string]string),
		"",
	)
}
