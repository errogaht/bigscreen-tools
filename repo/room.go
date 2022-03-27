package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/jackc/pgx/v4"
	"os"
)

type Room struct {
	Conn *pgx.Conn
}

func (repo *Room) DeleteAll() {
	_, err := repo.Conn.Exec(context.Background(), "DELETE FROM rooms;")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
}

func (repo *Room) Insert(rooms *[]bs.Room) {
	batch := &pgx.Batch{}
	for _, room := range *rooms {
		batch.Queue(
			"insert into rooms (id, created_at, participants, status, invite_code, visibility, room_type, version, size, environment, category, description, name, creator_profile) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
			room.Id, room.CreatedAt, room.Participants, room.Status, room.InviteCode, room.Visibility, room.RoomType, room.Version,
			room.Size, room.Environment, room.Category, room.Description, room.Name, room.CreatorProfile.Username,
		)
	}

	br := repo.Conn.SendBatch(context.Background(), batch)
	_, err := br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}
