package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/errogaht/bigscreen-tools/db"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

type OculusProfile struct {
	Conn *pgx.Conn
}

func (repo *OculusProfile) getMetadata() *db.TableMetadata {
	return &db.TableMetadata{
		Name: "oculus_profiles",
		Cols: []string{"id", "image_url", "small_image_url"},
		PK:   "id",
	}
}

func (r *OculusProfile) findBy(cond string, args ...interface{}) *[]bs.OculusProfile {
	md := r.getMetadata()
	sql := md.GetFindBySql(cond)
	rows, err := r.Conn.Query(context.Background(), sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var rowSlice []bs.OculusProfile

	for rows.Next() {
		var p bs.OculusProfile

		err := rows.Scan(&p.Id, &p.ImageURL, &p.SmallImageURL)
		if err != nil {
			log.Fatal(err)
		}
		rowSlice = append(rowSlice, p)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return &rowSlice
}

func (repo *OculusProfile) Upsert(profiles *[]bs.OculusProfile) {
	md := repo.getMetadata()
	batch := &pgx.Batch{}
	for _, p := range *profiles {
		batch.Queue(
			md.GetUpsertSql(),
			p.Id, p.ImageURL, p.SmallImageURL,
		)
	}

	br := repo.Conn.SendBatch(context.Background(), batch)
	defer br.Close()
	_, err := br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}

func (repo *OculusProfile) GetProfilesFrom(rooms *[]bs.Room) (profiles []bs.OculusProfile) {
	profilesSet := make(map[string]struct{})
	var p *bs.OculusProfile
	for i := range *rooms {
		p = &(*rooms)[i].CreatorProfile.OculusProfile
		if p.Id == "" {
			continue
		}
		if _, ok := profilesSet[p.Id]; ok {
			continue
		}
		profilesSet[p.Id] = struct{}{}
		profiles = append(profiles, *p)
	}
	return
}
