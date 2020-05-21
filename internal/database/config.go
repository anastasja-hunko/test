package database

type Config struct {
	DatabaseURL  string
	DatabaseName string
	UserColName  string
	DocColName   string
}

func NewConfig() *Config {
	return &Config{
		DatabaseURL:  "mongodb://localhost:27017/test_task",
		DatabaseName: "test_task",
		UserColName:  "users",
		DocColName:   "docs",
	}
}
