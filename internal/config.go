package internal

import (
	"github.com/urfave/cli"
)

func NewConfig(c *cli.Context) *Config {
	return &Config{
		SpaceDomain: c.String("host"),
		ProjectKey:  c.String("project"),
		APIKey:      c.String("key"),
		File:        c.String("file"),
		Lang: func() Lang {
			if c.String("lang") == En.Value() {
				return En
			} else {
				return Ja
			}
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
