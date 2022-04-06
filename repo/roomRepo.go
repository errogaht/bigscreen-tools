package repo

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/errogaht/bigscreen-tools/common"
	"github.com/errogaht/bigscreen-tools/db"
	"github.com/errogaht/bigscreen-tools/s"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"time"
)

type RoomRepo struct {
	conn               *pgx.Conn
	accountProfileRepo *AccountProfileRepo
	oculusProfilesRepo *OculusProfileRepo
	steamProfilesRepo  *SteamProfileRepo
	settingsRepo       *SettingsRepo
}

func NewRoomRepo(
	conn *pgx.Conn,
	accountProfileRepo *AccountProfileRepo,
	oculusProfilesRepo *OculusProfileRepo,
	steamProfilesRepo *SteamProfileRepo,
	settingsRepo *SettingsRepo,
) *RoomRepo {
	return &RoomRepo{
		conn:               conn,
		accountProfileRepo: accountProfileRepo,
		oculusProfilesRepo: oculusProfilesRepo,
		steamProfilesRepo:  steamProfilesRepo,
		settingsRepo:       settingsRepo,
	}
}

type Cols []string

func (repo *RoomRepo) DeleteAll() {
	_, err := repo.conn.Exec(context.Background(), "DELETE FROM rooms;")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
}

func (repo *RoomRepo) RefreshRoomsInDB(rooms *[]bs.Room) {
	oculusProfiles := bs.GetOculusProfilesFrom(rooms)
	repo.oculusProfilesRepo.Upsert(&oculusProfiles)
	common.LogM(fmt.Sprintf("%d oculusProfiles upsert", len(oculusProfiles)))

	steamProfiles := bs.GetSteamProfilesFrom(rooms)
	repo.steamProfilesRepo.Upsert(&steamProfiles)
	common.LogM(fmt.Sprintf("%d steamProfiles upsert", len(steamProfiles)))

	creatorProfiles := bs.GetCreatorProfilesFrom(rooms)
	repo.accountProfileRepo.Upsert(&creatorProfiles)

	repo.DeleteAll()
	repo.Insert(rooms)
	repo.settingsRepo.Upsert(&[]s.Settings{{Id: SETTING_ROOMS_LAST_UPDATED, Timestamp: time.Now()}})
	common.LogM(fmt.Sprintf("rooms refreshed, 30s. sleep..."))
}

func (repo *RoomRepo) getMetadata() *db.TableMetadata {
	return &db.TableMetadata{
		Name: "rooms",
		Cols: []string{"id", "created_at", "participants", "status", "invite_code", "visibility", "room_type", "version", "size", "environment", "category", "description", "name", "creator_profile"},
		PK:   "id",
	}
}

func (repo *RoomRepo) FindBy(cond string, args ...interface{}) *[]bs.Room {
	md := repo.getMetadata()
	sql := md.GetFindBySql(cond)
	rows, err := repo.conn.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}
	defer rows.Close()

	var rowSlice []bs.Room
	var accProfilesIds []interface{}
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
		rowSlice = append(rowSlice, r)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}
	if len(rowSlice) == 0 {
		null := make([]bs.Room, 0)
		return &null
	}

	accProfiles := repo.accountProfileRepo.findBy(fmt.Sprintf("username IN(%s)", md.Params2(accProfilesIds)), accProfilesIds...)
	creatorProfilesById := make(map[string]*bs.AccountProfile)
	for i := range *accProfiles {
		profile := &(*accProfiles)[i]
		creatorProfilesById[profile.Username] = profile
	}
	for i := range rowSlice {
		r := &rowSlice[i]
		if profile, ok := creatorProfilesById[r.CreatorProfileUsername]; ok {
			r.CreatorProfile = *profile
		}
	}
	return &rowSlice
}
func (repo *RoomRepo) Insert(rooms *[]bs.Room) {
	batch := &pgx.Batch{}
	md := repo.getMetadata()
	for _, room := range *rooms {
		batch.Queue(
			fmt.Sprintf(`insert into rooms (%s) values(%s)`, md.Comma(), md.Params()),
			room.Id, room.CreatedAt, room.Participants, room.Status, room.InviteCode, room.Visibility, room.RoomType, room.Version,
			room.Size, room.Environment, room.Category, room.Description, room.Name, room.CreatorProfile.Username,
		)
	}

	br := repo.conn.SendBatch(context.Background(), batch)
	defer br.Close()
	_, err := br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}
