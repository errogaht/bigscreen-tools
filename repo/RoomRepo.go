package repo

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/errogaht/bigscreen-tools/common"
	"github.com/errogaht/bigscreen-tools/s"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"time"
)

type RoomRepo struct {
	Repo
	accountProfileRepo *AccountProfileRepo
	oculusProfilesRepo *OculusProfileRepo
	steamProfilesRepo  *SteamProfileRepo
	roomUsersRepo      *RoomUsersRepo
	settingsRepo       *SettingsRepo
}

func NewRoomRepo(
	conn *pgx.Conn,
	accountProfileRepo *AccountProfileRepo,
	oculusProfilesRepo *OculusProfileRepo,
	steamProfilesRepo *SteamProfileRepo,
	roomUsersRepo *RoomUsersRepo,
	settingsRepo *SettingsRepo,
) *RoomRepo {
	return &RoomRepo{
		Repo: Repo{
			Conn: conn,
			TblMetadata: TableMetadata{
				Name: "rooms",
				Cols: []string{"id", "created_at", "participants", "status", "invite_code", "visibility", "room_type", "version", "size", "environment", "category", "description", "name", "creator_profile"},
				PK:   "id",
			},
		},
		accountProfileRepo: accountProfileRepo,
		oculusProfilesRepo: oculusProfilesRepo,
		steamProfilesRepo:  steamProfilesRepo,
		settingsRepo:       settingsRepo,
		roomUsersRepo:      roomUsersRepo,
	}
}

type Cols []string

func (repo *RoomRepo) RefreshRoomsInDB(rooms *[]bs.Room) {
	oculusProfiles := bs.GetOculusProfilesFrom(rooms)
	repo.oculusProfilesRepo.Upsert(&oculusProfiles)
	common.LogM(fmt.Sprintf("%d oculusProfiles upsert", len(oculusProfiles)))

	steamProfiles := bs.GetSteamProfilesFrom(rooms)
	repo.steamProfilesRepo.Upsert(&steamProfiles)
	common.LogM(fmt.Sprintf("%d steamProfiles upsert", len(steamProfiles)))

	creatorProfiles := bs.GetAccountProfilesFrom(rooms)
	repo.accountProfileRepo.Upsert(&creatorProfiles)

	roomUsers := bs.GetRoomUsersFrom(rooms)
	repo.roomUsersRepo.DeleteAll() //live list of room users
	repo.roomUsersRepo.Upsert(&roomUsers)
	common.LogM(fmt.Sprintf("%d roomUsers refreshed", len(roomUsers)))

	repo.DeleteAll() //live list of rooms
	repo.Insert(rooms)
	repo.settingsRepo.Upsert(&[]s.Settings{{Id: SETTING_ROOMS_LAST_UPDATED, Timestamp: time.Now()}})
	common.LogM(fmt.Sprintf("rooms refreshed, 30s. sleep..."))
}

func (repo *RoomRepo) FindBy(cond string, args ...interface{}) *[]bs.Room {
	sql := repo.TblMetadata.GetFindBySql(cond)
	rows, err := repo.Conn.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}
	defer rows.Close()

	var rooms []bs.Room
	var accProfilesIds []interface{}
	var roomIds []interface{}
	var creatorProfileUsername sql2.NullString
	for rows.Next() {
		var r bs.Room
		err := rows.Scan(&r.Id, &r.CreatedAt, &r.Participants, &r.Status, &r.InviteCode, &r.Visibility, &r.RoomType, &r.Version, &r.Size, &r.Environment, &r.Category, &r.Description, &r.Name, &creatorProfileUsername)
		if err != nil {
			fmt.Printf("%v\n", sql)
			fmt.Printf("%v\n", args)
			log.Fatal(err)
		}
		if creatorProfileUsername.Valid {
			r.CreatorProfileUsername = creatorProfileUsername.String
			accProfilesIds = append(accProfilesIds, r.CreatorProfileUsername)
		}
		roomIds = append(roomIds, r.RoomId)
		rooms = append(rooms, r)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}
	if len(rooms) == 0 {
		null := make([]bs.Room, 0)
		return &null
	}

	roomUsers := repo.roomUsersRepo.FindBy(fmt.Sprintf("room_id IN(%s)", repo.TblMetadata.SqlParamsFrom(roomIds)))

	//collect accProfiles from roomUsers to accProfilesIds
	for i := range *roomUsers {
		accProfilesIds = append(accProfilesIds, (*roomUsers)[i].AccountProfileId)
	}

	//get all required accProfiles from DB
	accProfiles := repo.accountProfileRepo.findBy(fmt.Sprintf("username IN(%s)", repo.TblMetadata.SqlParamsFrom(accProfilesIds)), accProfilesIds...)

	//map accProfiles by their ids
	accountProfilesById := make(map[string]*bs.AccountProfile)
	for i := range *accProfiles {
		profile := &(*accProfiles)[i]
		accountProfilesById[profile.Username] = profile
	}

	//map slices of roomUsers by roomId they belongs to
	roomUsersByRId := make(map[string][]bs.RoomUser)
	for i := range *roomUsers {
		ru2 := &(*roomUsers)[i]

		//attach accountProfile to roomUser
		if ap, ok := accountProfilesById[ru2.AccountProfileId]; ok {
			ru2.AccountProfile = *ap
		}

		roomUsersByRId[ru2.RoomId] = append(roomUsersByRId[ru2.RoomId], *ru2)
	}

	for i := range rooms {
		r := &rooms[i]

		//attach accProfile to room
		if profile, ok := accountProfilesById[r.CreatorProfileUsername]; ok {
			r.CreatorProfile = *profile
		}

		//attach roomUsers to room
		if roomUsersPack, ok := roomUsersByRId[r.Id]; ok {
			r.RemoteUsers = roomUsersPack
		}
	}
	return &rooms
}

func (repo *RoomRepo) Insert(rooms *[]bs.Room) {
	batch := &pgx.Batch{}
	for _, room := range *rooms {
		batch.Queue(
			fmt.Sprintf("insert into %s (%s) values(%s)", repo.TblMetadata.Name, repo.TblMetadata.Comma(), repo.TblMetadata.SqlParams()),
			room.Id, room.CreatedAt, room.Participants, room.Status, room.InviteCode, room.Visibility, room.RoomType, room.Version,
			room.Size, room.Environment, room.Category, room.Description, room.Name, room.CreatorProfile.Username,
		)
	}

	br := repo.Conn.SendBatch(context.Background(), batch)
	defer br.Close()
	_, err := br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}
