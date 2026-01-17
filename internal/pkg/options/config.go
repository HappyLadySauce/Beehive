package options

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/HappyLadySauce/Beehive/pkg/utils/homedir"
)

const configFlagName = "config"

var cfgFile string

func init() {
	pflag.StringVarP(&cfgFile, "config", "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support JSON, TOML, YAML, HCL, or Java properties formats.")
}

// addConfigFlag adds flags for a specific server to the specified FlagSet
// object.
func AddConfigFlag(basename string, fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(configFlagName))

	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(basename), "-", "_", -1))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	cobra.OnInitialize(func() {
		if cfgFile != "" {
			// Support ${ENV_VAR} expansion inside config files.
			// This enables passing config values via environment variables (e.g. from make).
			b, err := os.ReadFile(cfgFile)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", cfgFile, err)
				os.Exit(1)
			}

			expanded := os.ExpandEnv(string(b))
			ext := strings.TrimPrefix(filepath.Ext(cfgFile), ".")
			if ext != "" {
				viper.SetConfigType(ext)
			}
			if err := viper.ReadConfig(strings.NewReader(expanded)); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", cfgFile, err)
				os.Exit(1)
			}
			return
		} else {
			viper.AddConfigPath(".")

			if names := strings.Split(basename, "-"); len(names) > 1 {
				viper.AddConfigPath(filepath.Join(homedir.HomeDir(), "."+names[0]))
				viper.AddConfigPath(filepath.Join("/etc", names[0]))
			}

			viper.SetConfigName(basename)
		}

		// Try to read config file, but don't fail if it doesn't exist
		// Configuration can come from environment variables or defaults
		if err := viper.ReadInConfig(); err != nil {
			// Check if it's a "file not found" error
			// viper returns an error when config file is not found, but we want to allow that
			if strings.Contains(err.Error(), "Not Found") ||
				strings.Contains(err.Error(), "no such file") ||
				strings.Contains(err.Error(), "cannot find") {
				// Config file not found is OK, we'll use defaults and env vars
				// This allows the application to run without a config file
			} else {
				// Other errors (like parse errors) should still be reported
				_, _ = fmt.Fprintf(os.Stderr, "Warning: failed to read configuration file: %v\n", err)
				_, _ = fmt.Fprintf(os.Stderr, "Continuing with defaults and environment variables...\n")
			}
		}
	})
}
