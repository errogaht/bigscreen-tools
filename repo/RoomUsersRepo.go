package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/jackc/pgx/v4"
	"os"
)

type RoomUsersRepo struct {
	Repo
}

func NewRoomUsersRepo(
	conn *pgx.Conn,
) *RoomUsersRepo {
	return &RoomUsersRepo{
		Repo: Repo{
			Conn: conn,
			TblMetadata: TableMetadata{
				Name: "room_users",
				Cols: []string{"user_session_id", "seat_index", "created_at", "account_profile", "room_id", "version", "is_staff", "is_mod", "is_admin"},
				PK:   "user_session_id",
			},
		},
	}
}

// FindBy does not attach AccountProfiles for optimization purposes. it's done at one level above
func (repo *RoomUsersRepo) FindBy(cond string, args ...interface{}) *[]bs.RoomUser {
	var rowSlice []bs.RoomUser
	repo.Repo.FindBy(cond, &rowSlice, args...)

	return &rowSlice
}

func (repo *RoomUsersRepo) Upsert(roomUsers *[]bs.RoomUser) {
	batch := &pgx.Batch{}
	for _, p := range *roomUsers {
		batch.Queue(
			repo.TblMetadata.GetUpsertSql(),
			p.UserSessionId, p.SeatIndex, p.CreatedAt, p.AccountProfileId, p.RoomId, p.Version, p.IsStaff, p.IsMod, p.IsAdmin,
		)
	}

	br := repo.Conn.SendBatch(context.Background(), batch)
	defer br.Close()
	_, err := br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}
