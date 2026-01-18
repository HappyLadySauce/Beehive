package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

// DatabaseOptions 数据库配置选项
type DatabaseOptions struct {
	Host         string `json:"host" mapstructure:"host"`
	Port         int    `json:"port" mapstructure:"port"`
	User         string `json:"user" mapstructure:"user"`
	Password     string `json:"password" mapstructure:"password"`
	DBName       string `json:"dbname" mapstructure:"dbname"`
	SSLMode      string `json:"sslmode" mapstructure:"sslmode"`
	MaxIdleConns int    `json:"max-idle-conns" mapstructure:"max-idle-conns"`
	MaxOpenConns int    `json:"max-open-conns" mapstructure:"max-open-conns"`
}


// NewDatabaseOptions 创建数据库配置选项
func NewDatabaseOptions() *DatabaseOptions {
	return &DatabaseOptions{
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


// Validate validates the DatabaseOptions.
func (o *DatabaseOptions) Validate() []error {
	var errs []error

	if o.Host == "" {
		errs = append(errs, fmt.Errorf("database.host cannot be empty"))
	}

	if o.Port < 1 || o.Port > 65535 {
		errs = append(errs, fmt.Errorf("database.port must be between 1 and 65535"))
	}

	if o.User == "" {
		errs = append(errs, fmt.Errorf("database.user cannot be empty"))
	}

	if o.DBName == "" {
		errs = append(errs, fmt.Errorf("database.dbname cannot be empty"))
	}

	if o.MaxIdleConns < 0 {
		errs = append(errs, fmt.Errorf("database.max-idle-conns cannot be negative"))
	}

	if o.MaxOpenConns < 1 {
		errs = append(errs, fmt.Errorf("database.max-open-conns must be at least 1"))
	}

	return errs
}


// AddFlags adds the flags to the specified FlagSet.
func (o *DatabaseOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "database.host", o.Host, "Database host address.")
	fs.IntVar(&o.Port, "database.port", o.Port, "Database port number.")
	fs.StringVar(&o.User, "database.user", o.User, "Database user name.")
	fs.StringVar(&o.Password, "database.password", o.Password, "Database password.")
	fs.StringVar(&o.DBName, "database.dbname", o.DBName, "Database name.")
	fs.StringVar(&o.SSLMode, "database.sslmode", o.SSLMode, "Database SSL mode (disable, require, verify-ca, verify-full).")
	fs.IntVar(&o.MaxIdleConns, "database.max-idle-conns", o.MaxIdleConns, "Maximum number of idle database connections.")
	fs.IntVar(&o.MaxOpenConns, "database.max-open-conns", o.MaxOpenConns, "Maximum number of open database connections.")
}
