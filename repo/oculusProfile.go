package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/errogaht/bigscreen-tools/db"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"strings"
)

type OculusProfileRepo struct {
	conn *pgx.Conn
}

func NewOculusProfileRepo(
	conn *pgx.Conn,
) *OculusProfileRepo {
	return &OculusProfileRepo{
		conn: conn,
	}
}

func (repo *OculusProfileRepo) getMetadata() *db.TableMetadata {
	return &db.TableMetadata{
		Name: "oculus_profiles",
		Cols: []string{"id", "image_url", "small_image_url"},
		PK:   "id",
	}
}

func (repo *OculusProfileRepo) findBy(cond string, args ...interface{}) *[]bs.OculusProfile {
	var rowSlice []bs.OculusProfile
	if strings.Contains(cond, "IN") && len(args) == 0 {
		return &rowSlice
	}
	md := repo.getMetadata()
	sql := md.GetFindBySql(cond)
	rows, err := repo.conn.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p bs.OculusProfile

		err := rows.Scan(&p.Id, &p.ImageURL, &p.SmallImageURL)
		if err != nil {
			fmt.Printf("%v\n", sql)
			fmt.Printf("%v\n", args)
			log.Fatal(err)
		}
		rowSlice = append(rowSlice, p)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}

	return &rowSlice
}

func (repo *OculusProfileRepo) Upsert(profiles *[]bs.OculusProfile) {
	md := repo.getMetadata()
	batch := &pgx.Batch{}
	for _, p := range *profiles {
		batch.Queue(
			md.GetUpsertSql(),
			p.Id, p.ImageURL, p.SmallImageURL,
		)
	}

	br := repo.conn.SendBatch(context.Background(), batch)
	defer br.Close()
	_, err := br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}
