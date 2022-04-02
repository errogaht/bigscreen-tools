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

func GetCreatorProfilesFrom(rooms *[]Room) (profiles []AccountProfile) {
	profilesSet := make(map[string]struct{})
	for i := range *rooms {
		p := &(*rooms)[i].CreatorProfile
		if _, ok := profilesSet[p.Username]; ok {
			continue
		}
		profilesSet[p.Username] = struct{}{}
		profiles = append(profiles, *p)
	}
	return
}

func GetOculusProfilesFrom(rooms *[]Room) (profiles []OculusProfile) {
	profilesSet := make(map[string]struct{})
	var p *OculusProfile
	for i := range *rooms {
		p = &(*rooms)[i].CreatorProfile.OculusProfile
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

func GetSteamProfilesFrom(rooms *[]Room) (profiles []SteamProfile) {
	profilesSet := make(map[string]struct{})
	for i := range *rooms {
		p := &(*rooms)[i].CreatorProfile.SteamProfile
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
