package config

type Configurations struct {
	User UserConfigurations
	DefaultRegion string
	MFASerial     string
	RoleArn       string
	SessionName   string
}

type UserConfigurations struct{
	AccKeyId string
	SecAccKey string
}