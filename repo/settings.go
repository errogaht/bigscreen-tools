package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/s"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
)

type SettingsRepo struct {
	Conn *pgx.Conn
}

const SETTING_ROOMS_LAST_UPDATED = "rooms.last.updated"

func NewSettingsRepo(
	conn *pgx.Conn,
) *SettingsRepo {
	return &SettingsRepo{
		Conn: conn,
	}
}

func (repo *SettingsRepo) Find(id string) (r s.Settings) {
	var id2 string
	var timestamp time.Time
	var int2 int64
	var string2 string
	var bool2 bool

	err := repo.Conn.QueryRow(context.Background(), "select id, timestamp, int, string, bool from settings where id=$1", id).Scan(&id2, &timestamp, &int2, &string2, &bool2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}

	r.Id = id2
	r.Timestamp = timestamp
	r.Int = int2
	r.String = string2
	r.Bool = bool2

	return
}

func (repo *SettingsRepo) Upsert(settings *[]s.Settings) {
	batch := &pgx.Batch{}
	for _, p := range *settings {
		batch.Queue(
			"insert into settings "+
				"(id, timestamp, int, string, bool) "+
				"values($1, $2, $3, $4, $5) "+
				"on conflict (id) do update set "+
				"timestamp = excluded.timestamp, int = excluded.int, string = excluded.string, bool = excluded.bool",
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
