package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/s"
	"github.com/jackc/pgx/v4"
	"os"
)

type SettingsRepo struct {
	Repo
}

const SETTING_ROOMS_LAST_UPDATED = "rooms.last.updated"

func NewSettingsRepo(
	conn *pgx.Conn,
) *SettingsRepo {
	return &SettingsRepo{
		Repo: Repo{
			Conn: conn,
			TblMetadata: TableMetadata{
				Name: "settings",
				Cols: []string{"id", "timestamp", "int", "string", "bool"},
				PK:   "id",
			},
		},
	}
}

func (repo *SettingsRepo) Find(id string) (setting s.Settings) {
	repo.Repo.Find(id, &setting)
	return
}

func (repo *SettingsRepo) Upsert(settings *[]s.Settings) {
	batch := &pgx.Batch{}
	for _, p := range *settings {
		batch.Queue(
			repo.TblMetadata.GetUpsertSql(),
			p.Id, p.Timestamp, p.Int, p.String, p.Bool,
		)
	}

	br := repo.Conn.SendBatch(context.Background(), batch)
	defer br.Close()
	_, err := br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}
