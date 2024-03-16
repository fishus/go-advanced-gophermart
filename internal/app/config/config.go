package config

import "time"

type Config struct {
	runAddr      string        // runAddr адрес и порт запуска сервиса
	accrualAddr  string        // accrualAddr адрес системы расчёта начислений
	databaseURI  string        // databaseURI адрес подключения к базе данных
	jwtSecretKey string        // jwtSecretKey for JWT
	jwtExpires   time.Duration // jwtExpires Срок действия JWT
	logLevel     string        // logLevel Log level
}

const DecimalExponent = 5

// TODO Вынести сюда дефолтные значения

func InitConfig() Config {
	conf := NewConfig()
	conf = parseFlags(conf)
	conf = parseEnvs(conf)
	return conf
}

func NewConfig() Config {
	return Config{
		jwtExpires: 15 * time.Minute, // Вынести в дефолтные значения
	}
}

func (c Config) RunAddr() string {
	return c.runAddr
}

func (c Config) SetRunAddr(addr string) Config {
	c.runAddr = addr
	return c
}

func (c Config) AccrualAddr() string {
	return c.accrualAddr
}

func (c Config) SetAccrualAddr(addr string) Config {
	c.accrualAddr = addr
	return c
}

func (c Config) DatabaseURI() string {
	return c.databaseURI
}

func (c Config) SetDatabaseURI(uri string) Config {
	c.databaseURI = uri
	return c
}

func (c Config) JWTSecretKey() string {
	return c.jwtSecretKey
}

func (c Config) SetJWTSecretKey(key string) Config {
	c.jwtSecretKey = key
	return c
}

func (c Config) JWTExpires() time.Duration {
	return c.jwtExpires
}

func (c Config) SetJWTExpires(d time.Duration) Config {
	c.jwtExpires = d
	return c
}

func (c Config) LogLevel() string {
	return c.logLevel
}

func (c Config) SetLogLevel(level string) Config {
	c.logLevel = level
	return c
}
