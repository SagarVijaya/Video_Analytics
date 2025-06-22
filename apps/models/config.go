package models

type Config struct {
	Database Database
	Server   Server
	Redis    Redis
	Metrics  Metric
	Mail     MailDeatils
}

type Server struct {
	Port int
}

type Redis struct {
	Port string
}

type Metric struct {
	Port int
}

type Database struct {
	Host string
	Port int
	Name string
	User string
	Pass string
}

type MailDeatils struct {
	Host string
	Port string
	From string
	Pass string
	To   string
}
