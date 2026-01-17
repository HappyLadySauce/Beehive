package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// LogOptions contains configuration for log file output and rotation
type LogOptions struct {
	// LogFile is the path to the log file. If empty, logs will be written to stderr.
	LogFile string `json:"log-file" mapstructure:"log-file"`

	// MaxSize is the maximum size in megabytes of the log file before it gets rotated.
	// Default is 100 MB.
	MaxSize int `json:"log-max-size" mapstructure:"log-max-size"`

	// MaxBackups is the maximum number of old log files to retain.
	// Default is 3.
	MaxBackups int `json:"log-max-backups" mapstructure:"log-max-backups"`

	// MaxAge is the maximum number of days to retain old log files.
	// Default is 28 days.
	MaxAge int `json:"log-max-age" mapstructure:"log-max-age"`

	// Compress determines if the rotated log files should be compressed using gzip.
	// Default is true.
	Compress bool `json:"log-compress" mapstructure:"log-compress"`
}

func NewLogOptions() *LogOptions {
	return &LogOptions{
		MaxSize:    100, // 100 MB
		MaxBackups: 3,
		MaxAge:     28, // 28 days
		Compress:   true,
	}
}

func (l *LogOptions) Validate() []error {
	var errors []error

	// LogFile can be empty - if empty, logs will be written to stderr (as documented)
	// Only validate MaxSize, MaxBackups, MaxAge if LogFile is specified
	if l.LogFile != "" {
		if l.MaxSize <= 0 {
			errors = append(errors, fmt.Errorf("log-max-size must be greater than 0"))
		}
		if l.MaxBackups <= 0 {
			errors = append(errors, fmt.Errorf("log-max-backups must be greater than 0"))
		}
		if l.MaxAge <= 0 {
			errors = append(errors, fmt.Errorf("log-max-age must be greater than 0"))
		}
		// Compress is a boolean, both true and false are valid values
		// No validation needed for Compress
	}
	return errors
}

func (l *LogOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&l.LogFile, "log-file", l.LogFile, "If non-empty, write log files to this location")
	fs.IntVar(&l.MaxSize, "log-max-size", l.MaxSize, "Maximum size in megabytes of the log file before it gets rotated. It defaults to 100 megabytes.")
	fs.IntVar(&l.MaxBackups, "log-max-backups", l.MaxBackups, "Maximum number of old log files to retain. The default is 3.")
	fs.IntVar(&l.MaxAge, "log-max-age", l.MaxAge, "Maximum number of days to retain old log files. The default is 28 days.")
	fs.BoolVar(&l.Compress, "log-compress", l.Compress, "Compress determines if the rotated log files should be compressed using gzip. The default is true.")
}
