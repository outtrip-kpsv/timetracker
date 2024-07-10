package config

import (
  "fmt"
  "github.com/jessevdk/go-flags"
)

type OptionsSrv struct {
  Host   string `short:"h" long:"host" description:"хост" default:"localhost" env:"HOST"`
  Port   string `short:"p" long:"port" description:"порт" default:"3000" env:"PORT"`
  Log    string `long:"logger-create" description:"logger-create output" default:"debug" env:"LOG"`
  DbHost string `long:"dbhost" description:"the db server host" default:"localhost" env:"DB_HOST"`
  DbPort string `long:"dbport" description:"the db server port" default:"5432" env:"DB_PORT"`
  PgUser string `long:"pguser" description:"the db user" default:"user_postgres" env:"POSTGRES_USER"`
  PgPass string `long:"pgpass" description:"the db pass" default:"pass" env:"POSTGRES_PASSWORD"`
  DbName string `long:"dbname" description:"the db name" default:"test" env:"POSTGRES_DB"`
}

type ConfSrv struct {
  Options OptionsSrv
}

func InitConfServ() (*ConfSrv, error) {
  var conf ConfSrv
  var opts OptionsSrv
  parser := flags.NewParser(&opts, flags.Default)
  _, err := parser.Parse()
  if err != nil {
    return nil, err
  }
  conf.Options = opts
  return &conf, nil
}

func (o *OptionsSrv) DbString() string {
  return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", o.PgUser, o.PgPass, o.DbHost, o.DbPort, o.DbName)
}

func (o *OptionsSrv) ServStr() string {
  return fmt.Sprintf("%s:%s", o.Host, o.Port)
}
