package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/jackc/pgx/v4"
	"os"
)

type OculusProfile struct {
	Conn *pgx.Conn
}

func (repo *OculusProfile) InsertOrUpdate(profiles *[]*bs.OculusProfile) {
	batch := &pgx.Batch{}
	for _, p := range *profiles {
		batch.Queue(
			"insert into oculus_profiles (id, image_url, small_image_url) values($1, $2, $3) on conflict (id) do update set image_url = excluded.image_url, small_image_url = excluded.small_image_url",
			p.OculusId, p.OculusImageURL, p.OculusSmallImageURL,
		)
	}

	br := repo.Conn.SendBatch(context.Background(), batch)
	_, err := br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}

func (repo *OculusProfile) GetProfilesFrom(rooms *[]bs.Room) (profiles []*bs.OculusProfile) {
	profilesSet := make(map[string]struct{})
	for _, room := range *rooms {
		p := &room.CreatorProfile.OculusProfile
		if p.OculusId == "" {
			continue
		}
		if _, ok := profilesSet[p.OculusId]; ok {
			continue
		}
		profilesSet[p.OculusId] = struct{}{}
		profiles = append(profiles, p)
	}
	return
}
