package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/jackc/pgx/v4"
	"os"
)

type AccountProfile struct {
	Conn *pgx.Conn
}

func (repo *AccountProfile) InsertOrUpdate(profiles *[]bs.AccountProfile) {
	batch := &pgx.Batch{}
	for _, p := range *profiles {
		var steamId interface{}
		var oculusId interface{}
		if p.SteamProfile.Id == "" {
			steamId = nil
		} else {
			steamId = p.SteamProfile.Id
		}
		if p.OculusProfile.OculusId == "" {
			oculusId = nil
		} else {
			oculusId = p.OculusProfile.OculusId
		}
		batch.Queue(
			"insert into account_profiles "+
				"(username, created_at, is_verified, is_banned, is_staff, steam_profile_id, oculus_profile_id) "+
				"values($1, $2, $3, $4, $5, $6, $7) "+
				"on conflict (username) do update set "+
				"created_at = excluded.created_at, is_verified = excluded.is_verified, is_banned = excluded.is_banned, is_staff = excluded.is_staff, steam_profile_id = excluded.steam_profile_id, oculus_profile_id = excluded.oculus_profile_id",
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

func (repo *AccountProfile) GetCreatorProfilesFrom(rooms *[]bs.Room) (profiles []bs.AccountProfile) {
	profilesSet := make(map[string]struct{})
	for _, room := range *rooms {
		p := &room.CreatorProfile
		if _, ok := profilesSet[p.Username]; ok {
			continue
		}
		profilesSet[p.Username] = struct{}{}
		profiles = append(profiles, *p)
	}
	return
}
