// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"fmt"

	"github.com/elastic/beats/libbeat/cfgfile"
)

type TwitterbeatConfig struct {
	Period  *string `yaml:"period"`
	Twitter struct {
		Names          *[]string
		AccessKey      *string `yaml:"access_key"`
		AccessSecret   *string `yaml:"access_secret"`
		ConsumerKey    *string `yaml:"consumer_key"`
		ConsumerSecret *string `yaml:"consumer_secret"`
	}
}

type MandatoryConfigError struct {
	fieldname string
}

type yamlConfig struct {
	Twitterbeat TwitterbeatConfig
}

func NewTwitterbeatConfig() (*TwitterbeatConfig, error) {
	yaml := yamlConfig{}
	err := cfgfile.Read(&yaml, "")
	if err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	err = yaml.Twitterbeat.setDefaults()
	if err != nil {
		return nil, fmt.Errorf("Defaults could not be set: %v", err)
	}
	return &yaml.Twitterbeat, nil
}

func (c *TwitterbeatConfig) setDefaults() error {
	if c.Period == nil {
		*c.Period = "60s"
	}

	if c.Twitter.AccessKey == nil {
		return MandatoryConfigError{"AccessKey"}
	}

	if c.Twitter.AccessSecret == nil {
		return MandatoryConfigError{"AccessSecret"}
	}

	if c.Twitter.ConsumerKey == nil {
		return MandatoryConfigError{"ConsumerKey"}
	}

	if c.Twitter.ConsumerSecret == nil {
		return MandatoryConfigError{"ConsumerSecret"}
	}

	return nil
}

func (e MandatoryConfigError) Error() string {
	return fmt.Sprintf("Mandatory field \"%v\" was not set in config.", e.fieldname)
}
