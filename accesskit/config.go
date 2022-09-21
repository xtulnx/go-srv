package accesskit

// DbConfig 数据库
type DbConfig struct {
	DSN     string `json:"dsn" toml:"dsn"`
	Dialect string `json:"dialect" toml:"dialect"`
}
