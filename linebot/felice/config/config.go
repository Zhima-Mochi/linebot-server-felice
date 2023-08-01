package config

type Config struct {
	LineChannelSecret      string
	LineChannelToken       string
	OpenaiAPIKey           string
	CacheURL               string
	LinebotPort            string
	LineAdminUserIDList    []string
	LineCustomerUserIDList []string
}
