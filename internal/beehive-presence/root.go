package beehivePresence

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-presence/options"
)

func NewAPICommand(basename string, ctx context.Context) *cobra.Command {
	opts := options.NewOptions(basename)
	cmd := &cobra.Command{
		Use:   basename,
		Short: "BeehivePresence is a presence service for Beehive",
		Long:  "BeehivePresence is a presence service for Beehive",
		RunE: func(cmd *cobra.Command, args []string) error {
			// bind command line flags to viper (command line args override config file)
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			if err := viper.Unmarshal(opts); err != nil {
				return err
			}

			// initialize logs after flags are parsed and config is loaded
			logs.InitLogs()
			defer logs.FlushLogs()

			// setup log file rotation if log file is specified
			// This must be called after InitLogs() to ensure the log file setting takes effect
			if opts.Log.LogFile != "" {
				logWriter := &lumberjack.Logger{
					Filename:   opts.Log.LogFile,
					MaxSize:    opts.Log.MaxSize, // megabytes
					MaxBackups: opts.Log.MaxBackups,
					MaxAge:     opts.Log.MaxAge, // days
					Compress:   opts.Log.Compress,
				}
				klog.SetOutput(logWriter)
			}

			// validate options after flags & config are fully populated
			if errs := opts.Validate(); len(errs) != 0 {
				for _, err := range errs {
					fmt.Fprintln(os.Stderr, "Error:", err)
				}
				os.Exit(1)
			}
			return run(ctx, opts)
		},
	}

	nfs := opts.AddFlags(cmd.Flags())
	flag.SetUsageAndHelpFunc(cmd, *nfs, 80)

	return cmd
}

func run(ctx context.Context, opts *options.Options) error {

	return nil
}

