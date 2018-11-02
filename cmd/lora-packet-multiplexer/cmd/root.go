package cmd

import (
	"bytes"
	"io/ioutil"

	"github.com/brocaar/lora-packet-multiplexer/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var version string

var rootCmd = &cobra.Command{
	Use:   "lora-packet-multiplexer",
	Short: "LoRa Packet Multiplexer",
	Long: `LoRa Packet Multiplexer sends packet-forwarder data to multiple backends
	> documentation & support: https://www.loraserver.io/
	> source & copyright information: https://github.com/brocaar/lora-packet-multiplexer/`,
	RunE: run,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to configuration file (optional)")
	rootCmd.PersistentFlags().Int("log-level", 4, "debug=5, info=4, error=2, fatal=1, panic=0")

	viper.BindPFlag("general.log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	viper.SetDefault("packet_multiplexer.bind", "0.0.0.0:1700")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
}

// Execute executes the root command.
func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	config.Version = version

	if cfgFile != "" {
		b, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
		}

		viper.SetConfigType("toml")
		if err := viper.ReadConfig(bytes.NewBuffer(b)); err != nil {
			log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
		}
	} else {
		viper.SetConfigName("lora-packet-multiplexer")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/./config/lora-packet-multiplexer")
		viper.AddConfigPath("/etc/lora-packet-multiplexer")
		if err := viper.ReadInConfig(); err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				log.Warning("No configuration file found, using defaults.")
			default:
				log.WithError(err).Fatal("read configuration file error")
			}
		}
	}

	if err := viper.Unmarshal(&config.C); err != nil {
		log.WithError(err).Fatal("unmarshal config error")
	}

	log.SetLevel(log.Level(config.C.General.LogLevel))
}
