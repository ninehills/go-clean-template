package config

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/ilyakaznacheev/cleanenv"
)

//nolint:gochecknoglobals
var (
	cfg     *Config
	cfgLock = new(sync.RWMutex)
)

type (
	// Config -.
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Log   `yaml:"log"`
		MySQL `yaml:"mysql"` //nolint: tagliatelle
		Redis `yaml:"redis"`
	}

	// App -.
	App struct {
		Name      string `env:"APP_NAME"       env-required:"true" yaml:"name"`
		Debug     bool   `env:"APP_DEBUG"      yaml:"debug"`
		SuperUser string `env:"APP_SUPER_USER" env-required:"true" yaml:"superUser"`
		// Please changed when app first started.
		SuperPassword string `env:"APP_SUPER_PASSWORD" env-required:"true" yaml:"superPassword"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env:"HTTP_PORT" env-required:"true" yaml:"port"`
	}

	// Log -.
	Log struct {
		// debug, info, warn, error
		Level string `env:"LOG_LEVEL" env-required:"true" yaml:"level"`
		// output format, json/text
		Format string `env:"LOG_FORMAT" env-required:"true" yaml:"format"`
		// nocolor, default false
		NoColor bool `env:"LOG_NOCOLOR" yaml:"noColor"`
	}

	// MySQL -.
	MySQL struct {
		// https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DSN string `env:"MYSQL_DSN" env-required:"true"`
		// https://github.com/go-sql-driver/mysql#important-settings
		ConnMaxLifetime int `env:"MYSQL_CONN_MAX_LIFETIME" env-required:"true" yaml:"connMaxLifetime"`
		MaxOpenConns    int `env:"MYSQL_MAX_OPEN_CONNS"    env-required:"true" yaml:"maxOpenConns"`
		MaxIdleConns    int `env:"MYSQL_MAX_IDLE_CONNS"    env-required:"true" yaml:"maxIdleConns"`
	}

	// Redis -.
	Redis struct {
		// "redis://<user>:<pass>@localhost:6379/<db>"
		URL string `env:"REDIS_URL" env-required:"true"`
	}
)

func GetConfig() *Config {
	cfgLock.RLock()
	defer cfgLock.RUnlock()

	return cfg
}

// default cfg file: "./config/config.yml".
func LoadConfig(cfgFile string) (*Config, error) {
	c := &Config{}
	err := cleanenv.ReadConfig(cfgFile, c)
	if err != nil {
		return cfg, fmt.Errorf("config error: %w", err)
	}

	// 更新全局的cfg
	cfgLock.Lock()
	cfg = c
	cfgLock.Unlock()

	return cfg, nil
}

// Auto reload config.
// copy from https://github.com/spf13/viper/blob/v1.12.0/viper.go#L431
func Watcher(filename string, handler func(cfg *Config)) { //nolint: gocognit, cyclop
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		// we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
		configFile := filepath.Clean(filename)
		configDir, _ := filepath.Split(configFile)
		realConfigFile, _ := filepath.EvalSymlinks(filename)

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()

						return
					}
					currentConfigFile, _ := filepath.EvalSymlinks(filename)
					// we only care about the config file with the following cases:
					// 1 - if the config file was modified or created
					// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if (filepath.Clean(event.Name) == configFile &&
						event.Op&writeOrCreateMask != 0) ||
						(currentConfigFile != "" && currentConfigFile != realConfigFile) {
						realConfigFile = currentConfigFile
						log.Printf("Config file %s changed, reload it", event.Name)
						_, err := LoadConfig(filename)
						if err != nil {
							log.Fatalf("Load config error: %s, and crash!", err)
						}
						handler(cfg)
					} else if filepath.Clean(event.Name) == configFile &&
						event.Op&fsnotify.Remove != 0 {
						eventsWG.Done()

						return
					}

				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed
						log.Printf("watcher error: %v\n", err)
					}
					eventsWG.Done()

					return
				}
			}
		}()
		err = watcher.Add(configDir)
		if err != nil {
			panic(err)
		}
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}
