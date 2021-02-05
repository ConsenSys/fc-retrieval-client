package config

import (
	"flag"
	"github.com/spf13/pflag"
  "github.com/spf13/viper"
)

func NewConfig() *viper.Viper {
  conf := viper.New()
	conf.AutomaticEnv()
	defineFlags(conf)
	bindFlags(conf)	
	setValues(conf)
  return conf
}

func defineFlags(conf *viper.Viper) {
	flag.String("client-id", conf.GetString("CLIENT_ID"), "help message for client-id")
	flag.Int("ttl", conf.GetInt("ESTABLISHMENT_TTL"), "help message for ttl")
	flag.String("log-level", conf.GetString("LOG_LEVEL"), "help message for log-level")
	flag.String("log-target", conf.GetString("LOG_TARGET"), "help message for log-target")
}

func bindFlags(conf *viper.Viper) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	conf.BindPFlags(pflag.CommandLine)
}

func setValues(conf *viper.Viper) {
	conf.Set("CLIENT_ID", conf.GetString("client-id"))
	conf.Set("ESTABLISHMENT_TTL", conf.GetInt("ttl"))
	conf.Set("LOG_LEVEL", conf.GetString("log-level"))
	conf.Set("LOG_TARGET", conf.GetString("log-target"))
}
