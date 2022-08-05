package config

import (
	"fmt"
	"log"
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
		MySQL `yaml:"mysql"`
		Redis `yaml:"redis"`
	}

	// App -.
	App struct {
		Name      string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Debug     bool   `yaml:"debug" env:"APP_DEBUG"`
		SuperUser string `env-required:"true" yaml:"super_user" env:"APP_SUPER_USER"`
		// Please changed when app first started.
		SuperPassword string `env-required:"true" yaml:"super_password" env:"APP_SUPER_PASSWORD"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		// debug, info, warn, error
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
		// output format, json/text
		Format string `env-required:"true" yaml:"log_format"  env:"LOG_FORMAT"`
	}

	// MySQL -.
	MySQL struct {
		// https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DSN string `env-required:"true"                 env:"MYSQL_DSN"`
		// https://github.com/go-sql-driver/mysql#important-settings
		ConnMaxLifetime int `env-required:"true" yaml:"conn_max_lifetime" env:"MYSQL_CONN_MAX_LIFETIME"`
		MaxOpenConns    int `env-required:"true" yaml:"max_open_conns"     env:"MYSQL_MAX_OPEN_CONNS"`
		MaxIdleConns    int `env-required:"true" yaml:"max_idle_conns"     env:"MYSQL_MAX_IDLE_CONNS"`
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

// default cfg file: "./config/config.yml"
func LoadConfig(cfgFile string) (*Config, error) {
	c := &Config{}
	err := cleanenv.ReadConfig(cfgFile, c)
	if err != nil {
		return cfg, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(c)
	if err != nil {
		return cfg, err
	}
	// 更新全局的cfg
	cfgLock.Lock()
	cfg = c
	cfgLock.Unlock()
	return cfg, nil
}

// Auto reload config.
func ConfigWatcher(cfgFile string, handler func(cfg *Config)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("Config file %s changed, reload it", event.Name)
					_, err := LoadConfig(cfgFile)
					if err != nil {
						log.Fatalf("Load config error: %s, and crash!", err)
					}
					handler(cfg)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Fatalf("Config file watcher error: %s, and crash!", err)
			}
		}
	}()

	err = watcher.Add(cfgFile)
	if err != nil {
		log.Fatalf("Config file add watcher failed: %s, and crash!", cfgFile)
	}
	<-done
}
