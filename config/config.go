package config

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/ilyakaznacheev/cleanenv"
)

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
		Name      string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Debug     bool   `yaml:"debug" env:"APP_DEBUG"`
		SuperUser string `env-required:"true" yaml:"superUser" env:"APP_SUPER_USER"`
		// Please changed when app first started.
		SuperPassword string `env-required:"true" yaml:"superPassword" env:"APP_SUPER_PASSWORD"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		// debug, info, warn, error
		Level string `env-required:"true" yaml:"level"   env:"LOG_LEVEL"`
		// output format, json/text
		Format string `env-required:"true" yaml:"format"  env:"LOG_FORMAT"`
		// nocolor, default false
		NoColor bool `yaml:"noColor" env:"LOG_NOCOLOR"`
	}

	// MySQL -.
	MySQL struct {
		// https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DSN string `env-required:"true"                 env:"MYSQL_DSN"`
		// https://github.com/go-sql-driver/mysql#important-settings
		ConnMaxLifetime int `env-required:"true" yaml:"connMaxLifetime" env:"MYSQL_CONN_MAX_LIFETIME"`
		MaxOpenConns    int `env-required:"true" yaml:"maxOpenConns"     env:"MYSQL_MAX_OPEN_CONNS"`
		MaxIdleConns    int `env-required:"true" yaml:"maxIdleConns"     env:"MYSQL_MAX_IDLE_CONNS"`
	}

	// Redis -.
	Redis struct {
		// "redis://<user>:<pass>@localhost:6379/<db>"
		URL string `env-required:"true"                            env:"REDIS_URL"`
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

	err = cleanenv.ReadEnv(c)
	if err != nil {
		return cfg, fmt.Errorf("cleanv.ReadEnv faild : %w", err)
	}
	// 更新全局的cfg
	cfgLock.Lock()
	cfg = c
	cfgLock.Unlock()

	return cfg, nil
}

// Auto reload config.
// copy from https://github.com/spf13/viper/blob/v1.12.0/viper.go#L431
func ConfigWatcher(filename string, handler func(cfg *Config)) {
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
