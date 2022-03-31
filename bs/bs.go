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
