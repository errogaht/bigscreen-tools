package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/jackc/pgx/v4"
	"os"
)

type SteamProfileRepo struct {
	Repo
}

func NewSteamProfileRepo(
	conn *pgx.Conn,
) *SteamProfileRepo {
	return &SteamProfileRepo{
		Repo: Repo{
			Conn: conn,
			TblMetadata: TableMetadata{
				Name: "steam_profiles",
				Cols: []string{"id", "community_visibility_state", "profile_state", "persona_name", "profile_url", "avatar", "avatar_medium", "avatar_full", "avatar_hash", "persona_state", "real_name", "primary_clan_id", "created_at", "persona_state_flags", "loc_country_code"},
				PK:   "id",
			},
		},
	}
}

func (repo *SteamProfileRepo) FindBy(cond string, args ...interface{}) *[]bs.SteamProfile {
	var rowSlice []bs.SteamProfile
	repo.Repo.FindBy(cond, &rowSlice, args...)

	return &rowSlice
}

func (repo *SteamProfileRepo) Upsert(profiles *[]bs.SteamProfile) {
	batch := &pgx.Batch{}
	for _, p := range *profiles {
		batch.Queue(
			repo.TblMetadata.GetUpsertSql(),
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
