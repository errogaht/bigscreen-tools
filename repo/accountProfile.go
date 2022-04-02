package repo

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/errogaht/bigscreen-tools/db"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

type AccountProfile struct {
	Conn *pgx.Conn
}

func (repo *AccountProfile) getMetadata() *db.TableMetadata {
	return &db.TableMetadata{
		Name: "account_profiles",
		Cols: []string{"username", "created_at", "is_verified", "is_banned", "is_staff", "steam_profile_id", "oculus_profile_id"},
		PK:   "username",
	}
}

func (repo *AccountProfile) findBy(cond string, args ...interface{}) *[]bs.AccountProfile {
	md := repo.getMetadata()
	sql := md.GetFindBySql(cond)
	rows, err := repo.Conn.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}
	defer rows.Close()

	var rowSlice []bs.AccountProfile
	var steamProfilesIds []interface{}
	var oculusProfilesIds []interface{}
	for rows.Next() {
		var p bs.AccountProfile
		var steamProfileId sql2.NullString
		var oculusProfileId sql2.NullString
		err := rows.Scan(&p.Username, &p.CreatedAt, &p.IsVerified, &p.IsBanned, &p.IsStaff, &steamProfileId, &oculusProfileId)
		if err != nil {
			fmt.Printf("%v\n", sql)
			fmt.Printf("%v\n", args)
			log.Fatal(err)
		}
		if steamProfileId.Valid {
			steamProfilesIds = append(steamProfilesIds, steamProfileId.String)
			p.SteamProfileId = steamProfileId.String
		}

		if oculusProfileId.Valid {
			oculusProfilesIds = append(oculusProfilesIds, oculusProfileId.String)
			p.OculusProfileId = oculusProfileId.String
		}

		rowSlice = append(rowSlice, p)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}

	if len(rowSlice) == 0 {
		null := make([]bs.AccountProfile, 0)
		return &null
	}
	oculusProfilesRepo := OculusProfile{Conn: repo.Conn}
	oculusProfiles := oculusProfilesRepo.findBy(fmt.Sprintf("id IN(%s)", md.Params2(oculusProfilesIds)), oculusProfilesIds...)
	oculusProfilesById := make(map[string]*bs.OculusProfile)
	for i := range *oculusProfiles {
		profile := &(*oculusProfiles)[i]
		oculusProfilesById[profile.Id] = profile
	}

	steamProfilesRepo := SteamProfile{Conn: repo.Conn}
	steamProfiles := steamProfilesRepo.findBy(fmt.Sprintf("id IN(%s)", md.Params2(steamProfilesIds)), steamProfilesIds...)
	steamProfilesById := make(map[string]*bs.SteamProfile)
	for i := range *steamProfiles {
		profile := &(*steamProfiles)[i]
		steamProfilesById[profile.Id] = profile
	}
	for i := range rowSlice {
		accProfileRef := &rowSlice[i]
		if profile, ok := steamProfilesById[accProfileRef.SteamProfileId]; ok {
			accProfileRef.SteamProfile = *profile
		}
		if profile, ok := oculusProfilesById[accProfileRef.OculusProfileId]; ok {
			accProfileRef.OculusProfile = *profile
		}
	}

	return &rowSlice
}

func (repo *AccountProfile) Upsert(profiles *[]bs.AccountProfile) {
	md := repo.getMetadata()
	batch := &pgx.Batch{}
	for _, p := range *profiles {
		var steamId interface{}
		var oculusId interface{}
		if p.SteamProfile.Id == "" {
			steamId = nil
		} else {
			steamId = p.SteamProfile.Id
		}
		if p.OculusProfile.Id == "" {
			oculusId = nil
		} else {
			oculusId = p.OculusProfile.Id
		}
		batch.Queue(
			md.GetUpsertSql(),
			p.Username, p.CreatedAt, p.IsVerified, p.IsBanned, p.IsStaff, steamId, oculusId,
		)
	}

	br := repo.Conn.SendBatch(context.Background(), batch)
	defer br.Close()
	_, err := br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}
