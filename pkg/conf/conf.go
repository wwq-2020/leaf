package conf

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

const (
	ModeSegment                 = "segment"
	ModeSnowflake               = "snowflake"
	defaultListenAddr           = ":9001"
	defaultSnowflakeStep uint32 = 1000
)

var cfg conf

type conf struct {
	ListenAddr    string `toml:"listenaddr'`
	Mode          string `toml:"mode"`
	DBDatabase    string `toml:"dbdatabase"`
	DBUsername    string `toml:"dbusername"`
	DBPassword    string `toml:"dbpassword"`
	DBHost        string `toml:"dbhost"`
	DBPort        int    `toml:"dbport"`
	DBCharset     string `toml:"dbcharset"`
	LogLevel      string `toml:"loglevel"`
	SnowflakeStep uint32 `toml:"snowflakestep"`
}

func Parse(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrapf(err, "file", file)
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return errors.Wrapf(err, "file", file)
	}
	err = cfg.validateAndFill()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *conf) validateAndFill() error {
	if c.ListenAddr == "" {
		c.ListenAddr = defaultListenAddr
	}
	if c.Mode != ModeSegment && c.Mode != ModeSnowflake {
		c.Mode = ModeSnowflake
	}
	if c.Mode == ModeSegment &&
		GetDsn() == "" {
		return errors.New("dsn must be given when use segment mode")
	}
	return nil

}

func GetMode() string {
	return cfg.Mode
}

func GetListenAddr() string {
	return cfg.ListenAddr
}

func GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabase,
		cfg.DBCharset,
	)
}

func GetLogLevel() string {
	return cfg.LogLevel
}

func GetSnowflakeStep() uint32 {
	if cfg.SnowflakeStep == 0 {
		return defaultSnowflakeStep
	}
	return cfg.SnowflakeStep

}
