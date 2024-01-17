package config

import (
	"fmt"
)

type TCPServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (t TCPServer) Addr() string {
	return fmt.Sprintf("%s:%s", t.Host, t.Port)
}

type Logging struct {
	Level          string `yaml:"level"`
	Type           string `yaml:"type"`
	LogFileEnabled bool   `yaml:"logFileEnabled"`
	LogFilePath    string `yaml:"logFilePath"`
}

type SQL struct {
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Name        string `yaml:"name"`
	Port        string `yaml:"port"`
	MaxIdleConn int    `yaml:"maxIdleConn"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
}

func (s SQL) DatabaseUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		s.User, s.Password, s.Host, s.Port, s.Name)
}

func (s SQL) DataSourceName() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		s.User, s.Password, s.Host, s.Port, s.Name)
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
	// ConnDialTimeoutSec is Dial timeout for establishing new connections.
	// Default is 5 seconds.
	ConnDialTimeoutSec int `yaml:"connDialTimeoutSec"`
	// ReadTimeoutSec Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetReadDeadline calls completely.
	ReadTimeoutSec int `yaml:"readTimeoutSec"`
	// WriteTimeoutSec Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.  Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetWriteDeadline calls completely.
	WriteTimeoutSec int `yaml:"writeTimeoutSec"`
	// ConnPoolSize Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	ConnPoolSize int `yaml:"connPoolSize"`
	// ConnPoolTimeoutSec Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	ConnPoolTimeoutSec int `yaml:"connPoolTimeoutSec"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	// Default is 0. the idle connections are not closed by default.
	MinIdleConn int `yaml:"minIdleConn"`
	// Maximum number of idle connections.
	// Default is 0. the idle connections are not closed by default.
	MaxIdleConn int `yaml:"maxIdleConn"`
}

func (r Redis) Addr() string {
	return r.Host + ":" + r.Port
}

type Server struct {
	GRPC  TCPServer `yaml:"grpc"`
	HTTP  TCPServer `yaml:"http"`
	Log   Logging   `yaml:"log"`
	DB    SQL       `yaml:"db"`
	Redis Redis     `yaml:"redis"`
}
