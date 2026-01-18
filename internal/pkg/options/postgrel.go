package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

// DatabaseOptions 数据库配置选项
type PostgresqlOptions struct {
	Host         string `json:"host" mapstructure:"host"`
	Port         int    `json:"port" mapstructure:"port"`
	User         string `json:"user" mapstructure:"user"`
	Password     string `json:"password" mapstructure:"password"`
	DBName       string `json:"dbname" mapstructure:"dbname"`
	SSLMode      string `json:"sslmode" mapstructure:"sslmode"`
	MaxIdleConns int    `json:"max-idle-conns" mapstructure:"max-idle-conns"`
	MaxOpenConns int    `json:"max-open-conns" mapstructure:"max-open-conns"`
}


// NewPostgresqlOptions 创建数据库配置选项
func NewPostgresqlOptions() *PostgresqlOptions {
	return &PostgresqlOptions{
		Host:         "localhost",
		Port:         5432,
		User:         "postgres",
		Password:     "postgres",
		DBName:       "beehive",
		SSLMode:      "disable",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}
}


// Validate validates the PostgresqlOptions.
func (o *PostgresqlOptions) Validate() []error {
	var errs []error

	if o.Host == "" {
		errs = append(errs, fmt.Errorf("postgresql.host cannot be empty"))
	}

	if o.Port < 1 || o.Port > 65535 {
		errs = append(errs, fmt.Errorf("postgresql.port must be between 1 and 65535"))
	}

	if o.User == "" {
		errs = append(errs, fmt.Errorf("postgresql.user cannot be empty"))
	}

	if o.DBName == "" {
		errs = append(errs, fmt.Errorf("postgresql.dbname cannot be empty"))
	}

	if o.MaxIdleConns < 0 {
		errs = append(errs, fmt.Errorf("postgresql.max-idle-conns cannot be negative"))
	}

	if o.MaxOpenConns < 1 {
		errs = append(errs, fmt.Errorf("postgresql.max-open-conns must be at least 1"))
	}

	return errs
}


// AddFlags adds the flags to the specified PostgresqlOptions.
func (o *PostgresqlOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "postgresql.host", o.Host, "Database host address.")
	fs.IntVar(&o.Port, "postgresql.port", o.Port, "Database port number.")
	fs.StringVar(&o.User, "postgresql.user", o.User, "Database user name.")
	fs.StringVar(&o.Password, "postgresql.password", o.Password, "Database password.")
	fs.StringVar(&o.DBName, "postgresql.dbname", o.DBName, "Database name.")
	fs.StringVar(&o.SSLMode, "postgresql.sslmode", o.SSLMode, "Database SSL mode (disable, require, verify-ca, verify-full).")
	fs.IntVar(&o.MaxIdleConns, "postgresql.max-idle-conns", o.MaxIdleConns, "Maximum number of idle database connections.")
	fs.IntVar(&o.MaxOpenConns, "postgresql.max-open-conns", o.MaxOpenConns, "Maximum number of open database connections.")
}
