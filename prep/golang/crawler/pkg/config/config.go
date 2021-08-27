package config

import (
	"net/url"

	"github.com/spf13/pflag"
)

type Config struct {
	OutputDirectory string
	URL             *url.URL
}

func FromFlags() Config {
	cfg := Config{}
	pflag.StringVarP(&cfg.OutputDirectory, "output-directory", "o", "",
		"Specify the output directory for the local mirror.")
	pflag.Parse()
	return cfg
}
func (c *Config) SetURL(str string) error {
	u, err := url.Parse(str)
	if err != nil {
		return err
	}
	c.URL = u
	// fmt.Printf("%#v\n", c)
	return nil
}
