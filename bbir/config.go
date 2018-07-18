package bbir

import (
	"github.com/urfave/cli"
)

func NewConfig(c *cli.Context) *Config {
	return &Config{
		SpaceDomain: c.String("host"),
		ProjectKey:  c.String("project"),
		APIKey:      c.String("key"),
		File:        c.Args().First(),
		Lang: func() Lang {
			if c.String("lang") == Ja.Value() {
				return Ja
			}
			return En
		}(),
		Progress: c.Bool("progress"),
	}

}

type Config struct {
	SpaceDomain string
	ProjectKey  string
	APIKey      string
	File        string
	Lang        Lang
	Progress    bool
}

func (c *Config) HasFile() bool {
	return c.File != ""
}

type Lang string

func (l Lang) Value() string {
	return string(l)
}

const (
	Ja Lang = "ja"
	En Lang = "en"
)
