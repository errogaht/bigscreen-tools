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

type SteamProfile struct {
	Conn *pgx.Conn
}

func (repo *SteamProfile) getMetadata() *db.TableMetadata {
	return &db.TableMetadata{
		Name: "steam_profiles",
		Cols: []string{"id", "community_visibility_state", "profile_state", "persona_name", "profile_url", "avatar", "avatar_medium", "avatar_full", "avatar_hash", "persona_state", "real_name", "primary_clan_id", "created_at", "persona_state_flags", "loc_country_code"},
		PK:   "id",
	}
}

func (repo *SteamProfile) findBy(cond string, args ...interface{}) *[]bs.SteamProfile {
	var rowSlice []bs.SteamProfile
	if strings.Contains(cond, "IN") && len(args) == 0 {
		return &rowSlice
	}
	md := repo.getMetadata()
	sql := md.GetFindBySql(cond)
	rows, err := repo.Conn.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p bs.SteamProfile

		err := rows.Scan(&p.Id, &p.CommunityVisibilityState, &p.ProfileState, &p.PersonaName, &p.ProfileUrl, &p.Avatar, &p.AvatarMedium, &p.AvatarFull, &p.AvatarHash, &p.PersonaState, &p.RealName, &p.PrimaryClanId, &p.CreatedAt, &p.PersonaStateFlags, &p.LocCountryCode)
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

func (repo *SteamProfile) Upsert(profiles *[]bs.SteamProfile) {
	md := repo.getMetadata()
	batch := &pgx.Batch{}
	for _, p := range *profiles {
		batch.Queue(
			md.GetUpsertSql(),
			p.Id, p.CommunityVisibilityState, p.ProfileState, p.PersonaName, p.ProfileUrl, p.Avatar, p.AvatarMedium, p.AvatarFull, p.AvatarHash, p.PersonaState, p.RealName, p.PrimaryClanId, p.CreatedAt, p.PersonaStateFlags, p.LocCountryCode,
		)
	}

	br := repo.Conn.SendBatch(context.Background(), batch)
	defer br.Close()
	_, err := br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}
