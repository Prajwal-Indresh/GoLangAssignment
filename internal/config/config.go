package config

type Config struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBName     string `json:"db_name"`
	JWTSecret  string `json:"jwt_secret"`
}

func LoadConfig() (*Config, error) {
	return &Config{
		DBUser:     "root",
		DBPassword: "root",
		DBHost:     "localhost",
		DBPort:     "3306",
		DBName:     "students_db",
		JWTSecret:  "AssignGo",
	}, nil
}
