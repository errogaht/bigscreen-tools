package bs

//bs means Bigscreen, package for stuff related to bs

type Bigscreen struct {
	JWT          JWTToken
	Bearer       string
	HostAccounts string
	HostRealtime string
	Credentials  LoginCredentials
	DeviceInfo   string
	TgToken      string
}

func GetAccountProfilesFrom(rooms *[]Room) (profiles []AccountProfile) {
	profilesSet := make(map[string]struct{})
	for i := range *rooms {
		p := &(*rooms)[i].CreatorProfile

		appendAccountProfile(p, &profilesSet, &profiles)
		for i2 := range (*rooms)[i].RemoteUsers {
			p = &(*rooms)[i].RemoteUsers[i2].AccountProfile
			appendAccountProfile(p, &profilesSet, &profiles)
		}
	}
	return
}

func GetRoomUsersFrom(rooms *[]Room) (roomUsers []RoomUser) {
	ruSet := make(map[string]struct{})
	for i := range *rooms {
		for i2 := range (*rooms)[i].RemoteUsers {
			p := &(*rooms)[i].RemoteUsers[i2]

			if _, ok := ruSet[p.UserSessionId]; ok {
				return
			}
			ruSet[p.UserSessionId] = struct{}{}
			p.AccountProfileId = p.AccountProfile.Username
			p.RoomId = (*rooms)[i].RoomId
			roomUsers = append(roomUsers, *p)
		}
	}
	return
}

func GetOculusProfilesFrom(rooms *[]Room) (profiles []OculusProfile) {
	profilesSet := make(map[string]struct{})
	var p *OculusProfile
	for i := range *rooms {
		p = &(*rooms)[i].CreatorProfile.OculusProfile
		appendOculusProfile(p, &profilesSet, &profiles)
		for i2 := range (*rooms)[i].RemoteUsers {
			p = &(*rooms)[i].RemoteUsers[i2].AccountProfile.OculusProfile
			appendOculusProfile(p, &profilesSet, &profiles)
		}
	}
	return
}

func appendOculusProfile(p *OculusProfile, profilesSet *map[string]struct{}, profiles *[]OculusProfile) {
	if p.Id == "" {
		return
	}
	if _, ok := (*profilesSet)[p.Id]; ok {
		return
	}
	(*profilesSet)[p.Id] = struct{}{}
	*profiles = append(*profiles, *p)
}

func appendAccountProfile(p *AccountProfile, profilesSet *map[string]struct{}, profiles *[]AccountProfile) {
	if _, ok := (*profilesSet)[p.Username]; ok {
		return
	}
	(*profilesSet)[p.Username] = struct{}{}
	*profiles = append(*profiles, *p)
}

func appendSteamProfile(p *SteamProfile, profilesSet *map[string]struct{}, profiles *[]SteamProfile) {
	if p.Id == "" {
		return
	}
	if _, ok := (*profilesSet)[p.Id]; ok {
		return
	}
	(*profilesSet)[p.Id] = struct{}{}
	*profiles = append(*profiles, *p)
}

func GetSteamProfilesFrom(rooms *[]Room) (profiles []SteamProfile) {
	profilesSet := make(map[string]struct{})
	for i := range *rooms {
		p := &(*rooms)[i].CreatorProfile.SteamProfile
		appendSteamProfile(p, &profilesSet, &profiles)
		for i2 := range (*rooms)[i].RemoteUsers {
			p = &(*rooms)[i].RemoteUsers[i2].AccountProfile.SteamProfile
			appendSteamProfile(p, &profilesSet, &profiles)
		}
	}
	return
}
