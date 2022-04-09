package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/jackc/pgx/v4"
	"os"
)

type OculusProfileRepo struct {
	Repo
}

func NewOculusProfileRepo(
	conn *pgx.Conn,
) *OculusProfileRepo {
	return &OculusProfileRepo{
		Repo: Repo{
			Conn: conn,
			TblMetadata: TableMetadata{
				Name: "oculus_profiles",
				Cols: []string{"id", "image_url", "small_image_url"},
				PK:   "id",
			},
		},
	}
}

func (repo *OculusProfileRepo) findBy(cond string, args ...interface{}) *[]bs.OculusProfile {
	var rowSlice []bs.OculusProfile
	repo.Repo.FindBy(cond, &rowSlice, args...)

	return &rowSlice
}

func (repo *OculusProfileRepo) Upsert(profiles *[]bs.OculusProfile) {
	batch := &pgx.Batch{}
	for _, p := range *profiles {
		batch.Queue(
			repo.TblMetadata.GetUpsertSql(),
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
