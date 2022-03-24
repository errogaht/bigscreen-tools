package bs

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
)

type RoomCreator struct {
	IsMod   bool
	IsStaff bool
}

type AccountProfile struct {
	Username   string
	IsBanned   bool
	IsStaff    bool
	IsVerified bool
}

type RoomRemoteUsers struct {
	IsAdmin        bool
	IsMod          bool
	IsStaff        bool
	SeatIndex      uint8
	AccountProfile AccountProfile
}

type RoomCreatorProfileOculusProfile struct {
	OculusId       string
	OculusImageURL string
}
type RoomCreatorProfile struct {
	Username      string
	IsMod         bool
	IsStaff       bool
	IsVerified    bool
	OculusProfile RoomCreatorProfileOculusProfile
}
type Room struct {
	Name           string
	Creator        RoomCreator
	CreatorProfile RoomCreatorProfile
	Description    string
	Category       string
	Environment    string
	Visibility     string
	RoomType       string
	InviteCode     string
	CreatedAt      string
	RoomId         string
	Size           uint8
	Participants   uint8
	RemoteUsers    []RoomRemoteUsers
}

func (bsRef *Bigscreen) GetRooms() (listOfRooms []Room) {
	body, _ := bsRef.request(
		bsRef.HostAccounts+"/rooms/latest",
		"GET",
		make(map[string]string),
		"",
	)

	err := json.Unmarshal(body, &listOfRooms)
	if err != nil {
		log.Panic(err.Error())
	}

	return
}

func (bsRef *Bigscreen) PrintOnlineRooms() {
	listOfRooms := bsRef.GetRooms()

	sort.SliceStable(listOfRooms, func(i, j int) bool {
		return listOfRooms[i].Participants > listOfRooms[j].Participants
	})

	sort.SliceStable(listOfRooms, func(i, j int) bool {
		return listOfRooms[i].Category < listOfRooms[j].Category
	})

	for i, room := range listOfRooms {
		i++
		fmt.Printf(
			"%2d. %6v - %2d/%2d - %-30s | U:%-20s %s %s\n",
			i,
			room.Category,
			room.Participants,
			room.Size,
			room.Name,
			room.CreatorProfile.Username,
			room.Description,
			room.RoomId,
		)
	}
}

func (bsRef *Bigscreen) Participants(roomId string) (room Room) {
	bsRef.verify()
	body, _ := bsRef.request(
		bsRef.HostAccounts+"/room/"+roomId,
		"GET",
		make(map[string]string),
		"",
	)

	err := json.Unmarshal(body, &room)
	if err != nil {
		log.Panic(err.Error())
	}

	return
}

func (bsRef *Bigscreen) LeaveRoom() {
	bsRef.verify()
	bsRef.request(
		bsRef.HostAccounts+"/leave_room",
		"GET",
		make(map[string]string),
		"",
	)
}
