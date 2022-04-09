package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/randallmlough/pgxscan"
	"log"
	"os"
	"strings"
)

type Repo struct {
	Conn        *pgx.Conn
	TblMetadata TableMetadata
}

func (repo *Repo) DeleteAll() {
	_, err := repo.Conn.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s;", repo.TblMetadata.Name))
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func (repo *Repo) FindBy(cond string, result interface{}, args ...interface{}) {
	if strings.Contains(cond, "IN") && len(args) == 0 {
		return
	}
	sql := repo.TblMetadata.GetFindBySql(cond)
	rows, err := repo.Conn.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}
	defer rows.Close()

	if err := pgxscan.NewScanner(rows).Scan(result); err != nil {
		if err != nil {
			if err.Error() == "no rows in result set" {
				return
			}
			fmt.Printf("%v\n", sql)
			fmt.Printf("%v\n", args)
			log.Fatal(err)
		}
	}
}

func (repo *Repo) Find(id string, out interface{}) {
	sql := repo.TblMetadata.GetFindBySql(fmt.Sprintf("%s = $1", repo.TblMetadata.PK))
	rows, _ := repo.Conn.Query(context.Background(), sql, id)

	if err := pgxscan.NewScanner(rows).Scan(out); err != nil {
		if err != nil {
			fmt.Printf("%v\n", sql)
			fmt.Printf("%v\n", id)
			log.Fatal(err)
		}
	}

	return
}
