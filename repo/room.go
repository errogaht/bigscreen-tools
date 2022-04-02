package repo

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/errogaht/bigscreen-tools/db"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

type Room struct {
	Conn *pgx.Conn
}
type Cols []string

func (repo *Room) DeleteAll() {
	_, err := repo.Conn.Exec(context.Background(), "DELETE FROM rooms;")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
}

func (repo *Room) getMetadata() *db.TableMetadata {
	return &db.TableMetadata{
		Name: "rooms",
		Cols: []string{"id", "created_at", "participants", "status", "invite_code", "visibility", "room_type", "version", "size", "environment", "category", "description", "name", "creator_profile"},
		PK:   "id",
	}
}

func (repo *Room) FindBy(cond string, args ...interface{}) *[]bs.Room {
	md := repo.getMetadata()
	sql := md.GetFindBySql(cond)
	rows, err := repo.Conn.Query(context.Background(), sql, args...)
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
	accountProfilesRepo := AccountProfile{Conn: repo.Conn}

	accProfiles := accountProfilesRepo.findBy(fmt.Sprintf("username IN(%s)", md.Params2(accProfilesIds)), accProfilesIds...)
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
func (repo *Room) Insert(rooms *[]bs.Room) {
	batch := &pgx.Batch{}
	md := repo.getMetadata()
	for _, room := range *rooms {
		batch.Queue(
			fmt.Sprintf(`insert into rooms (%s) values(%s)`, md.Comma(), md.Params()),
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
