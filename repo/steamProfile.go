package repo

import (
	"context"
	"fmt"
	"github.com/errogaht/bigscreen-tools/bs"
	"github.com/jackc/pgx/v4"
	"os"
)

type SteamProfile struct {
	Conn *pgx.Conn
}

func (repo *SteamProfile) InsertOrUpdate(profiles *[]bs.SteamProfile) {
	batch := &pgx.Batch{}
	for _, p := range *profiles {
		batch.Queue(
			"insert into steam_profiles "+
				"(id, community_visibility_state, profile_state, persona_name, profile_url, avatar, avatar_medium, avatar_full, avatar_hash, persona_state, real_name, primary_clan_id, created_at, persona_state_flags, loc_country_code) "+
				"values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) "+
				"on conflict (id) do update set "+
				"community_visibility_state = excluded.community_visibility_state, profile_state = excluded.profile_state, persona_name = excluded.persona_name, profile_url = excluded.profile_url, avatar = excluded.avatar, avatar_medium = excluded.avatar_medium, avatar_full = excluded.avatar_full, avatar_hash = excluded.avatar_hash, persona_state = excluded.persona_state, real_name = excluded.real_name, primary_clan_id = excluded.primary_clan_id, created_at = excluded.created_at, persona_state_flags = excluded.persona_state_flags, loc_country_code = excluded.loc_country_code",
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

func (repo *SteamProfile) GetProfilesFrom(rooms *[]bs.Room) (profiles []bs.SteamProfile) {
	profilesSet := make(map[string]struct{})
	for _, room := range *rooms {
		p := &room.CreatorProfile.SteamProfile
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
