package cmd

import (
	"github.com/pallavJha/chaakoo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile    string
	verboseLog bool
	dryRun     bool

	rootCmd = &cobra.Command{
		Use:   "chaakoo",
		Short: "chaakoo converts the 2D grids or matrix into TMUX windows and panes",
		Run: func(cmd *cobra.Command, args []string) {
			var config *chaakoo.Config
			if err := viper.Unmarshal(config); err != nil {
				// TODO: add helpful example for a config
				log.Fatal().Err(err).Msg("cannot unmarshal the config")
			}
			if err := config.Validate(); err != nil {
				log.Fatal().Err(err).Msg("validation errors found in the config")
			}
			if err := config.Parse(); err != nil {
				log.Fatal().Err(err).Msg("cannot parse the grid for a window")
			}
			config.DryRun = dryRun
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is CWD/chaakoo.yaml)")
	if err := rootCmd.MarkPersistentFlagRequired("config"); err != nil {
		log.Fatal().Err(err).Msg("cannot set the flag config required")
	}
	rootCmd.PersistentFlags().BoolVar(&verboseLog, "verbose", false, "pass true to enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&verboseLog, "dry-run", false, "if true then commands will not be executed")
}

func initConfig() {
	reconfigureLogger()
	readConfig()
}

func readConfig() {
	if cfgFile != "" {
		log.Debug().Msgf("using %s", cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		log.Debug().Msg("config file is not provided, trying to read chaakoo.yaml from the working directory")
		workingDirectory, err := os.Getwd()
		if err != nil {
			log.Fatal().Err(err).Msg("cannot get the current working directory")
		}
		viper.AddConfigPath(workingDirectory)
		viper.SetConfigType("yaml")
		viper.SetConfigName("chaakoo")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msgf("cannot read the config file: %s", viper.ConfigFileUsed())
	}
	log.Debug().Msgf("using config file: %s", viper.ConfigFileUsed())
}

func reconfigureLogger() {
	log.Logger = zerolog.New(&zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Caller().Logger()
	if verboseLog {
		log.Debug().Msgf("setting global log level to DEBUG as verbose log is enabled")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}