package config

import (
	"testing"
)

var Config config

type config struct {
	runAddr     string // runAddr адрес и порт запуска сервиса
	accrualAddr string // accrualAddr адрес системы расчёта начислений
	databaseURI string // databaseURI адрес подключения к базе данных
}

func init() {
	if testing.Testing() {
		return
	}
	Config = initConfig()
}

func initConfig() config {
	conf := newConfig()
	conf = parseFlags(conf)
	conf = parseEnvs(conf)
	return conf
}

func newConfig() config {
	return config{}
}

func (c config) RunAddr() string {
	return c.runAddr
}

func (c config) SetRunAddr(addr string) config {
	c.runAddr = addr
	return c
}

func (c config) AccrualAddr() string {
	return c.accrualAddr
}

func (c config) SetAccrualAddr(addr string) config {
	c.accrualAddr = addr
	return c
}

func (c config) DatabaseURI() string {
	return c.databaseURI
}

func (c config) SetDatabaseURI(uri string) config {
	c.databaseURI = uri
	return c
}
