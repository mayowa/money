package internal

import (
	"path/filepath"

	"github.com/jinzhu/configor"
)

// Config ...
type Config struct {
	Addr  string `default:"localhost:4000"`
	Debug bool   `default:"true"`

	Options struct {
		Templates   bool `default:"true"`
		Sessions    bool
		Timeout     int
		LogRequests bool `default:"true"`
	}

	Session struct {
		Name     string `default:"session"`
		DSN      string
		Key      string `default:"chidinmaisafinegirl"`
		Duration int    `default:"24"`
		// 0: cookies, 1: database, 2: redis
		Type int
	}

	Folders struct {
		Public     string `default:"./public"`
		Resources  string `default:"./ui"`
		Templates  string `default:"./ui/templates"`
		Tmp        string `default:"./tmp"`
		Data       string `default:"./data"`
		Migrations string `default:"./migrations"`
	}

	URLS struct {
		Public string `default:"/public"`
	}

	DB struct {
		User     string `default:"sysdba"`
		Password string `default:"masterkey"`
		Host     string `default:"localhost"`
		Port     int    `default:"5432"`
	}
}

// ReadConfigOptions options for ReadConfig
type ReadConfigOptions struct {
	// Name the name of the config file
	Name string
	// Folder where the config file is located
	Folder string
}

// ReadConfig read config file set @options to use default options
func ReadConfig(options *ReadConfigOptions) (*Config, error) {

	if options == nil {
		options = &ReadConfigOptions{}
	}
	if options.Name == "" {
		options.Name = "config.yml"
	}

	cfgFile := filepath.Join(options.Folder, options.Name)
	cfg := new(Config)
	if err := configor.Load(cfg, cfgFile); err != nil {
		return nil, err
	}

	return cfg, nil
}
