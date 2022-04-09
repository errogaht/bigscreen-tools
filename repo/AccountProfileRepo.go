package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/jackc/pgx/v4"
	"github.com/randallmlough/pgxscan"
	"log"
	"os"
)

type AccountProfileRepo struct {
	Repo
	oculusProfileRepo *OculusProfileRepo
	steamProfileRepo  *SteamProfileRepo
}

func NewAccountProfileRepo(
	conn *pgx.Conn,
	oculusProfileRepo *OculusProfileRepo,
	steamProfileRepo *SteamProfileRepo,
) *AccountProfileRepo {
	return &AccountProfileRepo{
		Repo: Repo{
			Conn: conn,
			TblMetadata: TableMetadata{
				Name: "account_profiles",
				Cols: []string{"username", "created_at", "is_verified", "is_banned", "is_staff", "steam_profile_id", "oculus_profile_id"},
				PK:   "username",
			},
		},
		oculusProfileRepo: oculusProfileRepo,
		steamProfileRepo:  steamProfileRepo,
	}
}

func (repo *AccountProfileRepo) findBy(cond string, args ...interface{}) *[]bs.AccountProfile {
	sql := repo.TblMetadata.GetFindBySql(cond)
	rows, err := repo.Conn.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Printf("%v\n", sql)
		fmt.Printf("%v\n", args)
		log.Fatal(err)
	}
	defer rows.Close()

	var rowSlice []bs.AccountProfile

	if err := pgxscan.NewScanner(rows).Scan(&rowSlice); err != nil {
		if err != nil {
			fmt.Printf("%v\n", sql)
			fmt.Printf("%v\n", args)
			log.Fatal(err)
		}
	}

	if len(rowSlice) == 0 {
		return nil
	}

	var steamProfilesIds []interface{}
	var oculusProfilesIds []interface{}
	for i := range rowSlice {
		p := &rowSlice[i]
		if p.SteamProfileId.Valid {
			steamProfilesIds = append(steamProfilesIds, p.SteamProfileId.String)
		}

		if p.OculusProfileId.Valid {
			oculusProfilesIds = append(oculusProfilesIds, p.OculusProfileId.String)
		}
	}
	oculusProfiles := repo.oculusProfileRepo.findBy(fmt.Sprintf("id IN(%s)", repo.TblMetadata.SqlParamsFrom(oculusProfilesIds)), oculusProfilesIds...)
	oculusProfilesById := make(map[string]*bs.OculusProfile)
	for i := range *oculusProfiles {
		profile := &(*oculusProfiles)[i]
		oculusProfilesById[profile.Id] = profile
	}

	steamProfiles := repo.steamProfileRepo.FindBy(fmt.Sprintf("id IN(%s)", repo.TblMetadata.SqlParamsFrom(steamProfilesIds)), steamProfilesIds...)
	steamProfilesById := make(map[string]*bs.SteamProfile)
	for i := range *steamProfiles {
		profile := &(*steamProfiles)[i]
		steamProfilesById[profile.Id] = profile
	}
	for i := range rowSlice {
		accProfileRef := &rowSlice[i]
		if profile, ok := steamProfilesById[accProfileRef.SteamProfileId.String]; ok {
			accProfileRef.SteamProfile = *profile
		}
		if profile, ok := oculusProfilesById[accProfileRef.OculusProfileId.String]; ok {
			accProfileRef.OculusProfile = *profile
		}
	}

	return &rowSlice
}

func (repo *AccountProfileRepo) Upsert(profiles *[]bs.AccountProfile) {
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
			repo.TblMetadata.GetUpsertSql(),
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
