package app

import (
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FlagsTestSuite struct {
	suite.Suite
	osArgs    []string
	osEnviron map[string]string
}

func (suite *FlagsTestSuite) SetupSuite() {
	// Flags
	suite.osArgs = make([]string, 0)
	suite.osArgs = append(suite.osArgs, os.Args...)

	// ENV
	suite.osEnviron = make(map[string]string)
	for _, e := range []string{
		"RUN_ADDRESS",
		"ACCRUAL_SYSTEM_ADDRESS",
		"DATABASE_URI",
		"JWT_SECRET_KEY",
		"LOG_LEVEL",
	} {
		suite.osEnviron[e] = os.Getenv(e)
	}
}

func (suite *FlagsTestSuite) TearDownSuite() {
	// Flags
	os.Args = make([]string, 0)
	os.Args = append(os.Args, suite.osArgs...)

	// ENV
	for k, v := range suite.osEnviron {
		if v != "" {
			_ = os.Setenv(k, v)
		}
	}
}

func (suite *FlagsTestSuite) SetupSubTest() {
	// Flags
	os.Args = make([]string, 0)
	os.Args = append(os.Args, suite.osArgs[0])

	// Prepare a new default flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Clear ENV
	for k := range suite.osEnviron {
		_ = os.Unsetenv(k)
	}
}

func (suite *FlagsTestSuite) TestParseFlags() {
	testCases := []struct {
		name string
		args []string
		want map[string]interface{}
	}{
		{
			name: "Positive case: Default values",
			args: nil,
			want: map[string]interface{}{
				"runAddr":      "localhost:8080",
				"accrualAddr":  "localhost:8081",
				"databaseURI":  "",
				"jwtSecretKey": "MySecretKey",
				"logLevel":     "debug",
			},
		},
		{
			name: "Positive case: Set flag -a",
			args: []string{"-a=example.com:8180"},
			want: map[string]interface{}{"runAddr": "example.com:8180"},
		},
		{
			name: "Positive case: Set flag -r",
			args: []string{"-r=example.com:8181"},
			want: map[string]interface{}{"accrualAddr": "example.com:8181"},
		},
		{
			name: "Positive case: Set flag -d",
			args: []string{"-d=postgres://username:password@localhost:5432/database_name"},
			want: map[string]interface{}{"databaseURI": "postgres://username:password@localhost:5432/database_name"},
		},
		{
			name: "Positive case: Set flag -sk",
			args: []string{"-sk=Secret1"},
			want: map[string]interface{}{"jwtSecretKey": "Secret1"},
		},
		{
			name: "Positive case: Set flag -ll",
			args: []string{"-ll=fatal"},
			want: map[string]interface{}{"logLevel": "fatal"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if len(tc.args) > 0 {
				os.Args = append(os.Args, tc.args...)
			}

			config := newConfig()
			config = parseFlags(config)

			configFields := reflect.ValueOf(config)

			for k, want := range tc.want {
				field := configFields.FieldByName(k)
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'. Expected: '%+v'. Actual: '%+v'.", k, reflect.ValueOf(want), reflect.ValueOf(field))
			}
		})
	}
}

func (suite *FlagsTestSuite) TestParseEnvs() {
	testCases := []struct {
		name string
		envs []string
		want map[string]interface{}
	}{
		{
			name: "Positive case: Default values",
			envs: nil,
			want: map[string]interface{}{
				"runAddr":      "",
				"accrualAddr":  "",
				"databaseURI":  "",
				"jwtSecretKey": "",
				"logLevel":     "",
			},
		},
		{
			name: "Positive case: Set env RUN_ADDRESS",
			envs: []string{"RUN_ADDRESS=example.com:8180"},
			want: map[string]interface{}{"runAddr": "example.com:8180"},
		},
		{
			name: "Positive case: Set env ACCRUAL_SYSTEM_ADDRESS",
			envs: []string{"ACCRUAL_SYSTEM_ADDRESS=example.com:8181"},
			want: map[string]interface{}{"accrualAddr": "example.com:8181"},
		},
		{
			name: "Positive case: Set env DATABASE_URI",
			envs: []string{"DATABASE_URI=postgres://username:password@localhost:5432/database_name"},
			want: map[string]interface{}{"databaseURI": "postgres://username:password@localhost:5432/database_name"},
		},
		{
			name: "Positive case: Set env JWT_SECRET_KEY",
			envs: []string{"JWT_SECRET_KEY=Secret2"},
			want: map[string]interface{}{"jwtSecretKey": "Secret2"},
		},
		{
			name: "Positive case: Set env LOG_LEVEL",
			envs: []string{"LOG_LEVEL=fatal"},
			want: map[string]interface{}{"logLevel": "fatal"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if len(tc.envs) > 0 {
				for _, v := range tc.envs {
					e := strings.Split(v, "=")

					suite.Require().GreaterOrEqual(len(e), 2)

					if len(e) > 2 {
						e[1] = strings.Join(e[1:], "=")
						e = e[:2]
					}

					err := os.Setenv(e[0], e[1])
					suite.Require().NoError(err)
				}
			}

			config := newConfig()
			config = parseEnvs(config)

			configFields := reflect.ValueOf(config)

			for k, want := range tc.want {
				field := configFields.FieldByName(k)
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'. Expected: '%+v'. Actual: '%+v'.", k, reflect.ValueOf(want), reflect.ValueOf(field))
			}
		})
	}
}

func (suite *FlagsTestSuite) TestLoadConfig() {
	testCases := []struct {
		name string
		args []string
		envs []string
		want map[string]interface{}
	}{
		{
			name: "Positive case: Default values",
			args: nil,
			envs: nil,
			want: map[string]interface{}{
				"runAddr":      "localhost:8080",
				"accrualAddr":  "localhost:8081",
				"databaseURI":  "",
				"jwtSecretKey": "MySecretKey",
				"logLevel":     "debug",
			},
		},
		{
			name: "Positive case: Set flag -a and env RUN_ADDRESS",
			args: []string{"-a=aaa.com:1111"},
			envs: []string{"RUN_ADDRESS=bbb.com:2222"},
			want: map[string]interface{}{"runAddr": "bbb.com:2222"},
		},
		{
			name: "Positive case: Set flag -a only",
			args: []string{"-a=aaa.com:2222"},
			envs: nil,
			want: map[string]interface{}{"runAddr": "aaa.com:2222"},
		},
		{
			name: "Positive case: Set env RUN_ADDRESS only",
			args: nil,
			envs: []string{"RUN_ADDRESS=bbb.com:1111"},
			want: map[string]interface{}{"runAddr": "bbb.com:1111"},
		},
		{
			name: "Positive case: Set flag -r and env ACCRUAL_SYSTEM_ADDRESS",
			args: []string{"-r=aaa.com:3333"},
			envs: []string{"ACCRUAL_SYSTEM_ADDRESS=bbb.com:4444"},
			want: map[string]interface{}{"accrualAddr": "bbb.com:4444"},
		},
		{
			name: "Positive case: Set flag -r only",
			args: []string{"-r=aaa.com:4444"},
			envs: nil,
			want: map[string]interface{}{"accrualAddr": "aaa.com:4444"},
		},
		{
			name: "Positive case: Set env ACCRUAL_SYSTEM_ADDRESS only",
			args: nil,
			envs: []string{"ACCRUAL_SYSTEM_ADDRESS=bbb.com:3333"},
			want: map[string]interface{}{"accrualAddr": "bbb.com:3333"},
		},
		{
			name: "Positive case: Set flag -d and env DATABASE_URI",
			args: []string{"-d=postgres://username1:password1@localhost:5432/database_name1"},
			envs: []string{"DATABASE_URI=postgres://username2:password2@localhost:5432/database_name2"},
			want: map[string]interface{}{"databaseURI": "postgres://username2:password2@localhost:5432/database_name2"},
		},
		{
			name: "Positive case: Set flag -d only",
			args: []string{"-d=postgres://username1:password1@localhost:5432/database_name1"},
			envs: nil,
			want: map[string]interface{}{"databaseURI": "postgres://username1:password1@localhost:5432/database_name1"},
		},
		{
			name: "Positive case: Set env DATABASE_URI only",
			args: nil,
			envs: []string{"DATABASE_URI=postgres://username2:password2@localhost:5432/database_name2"},
			want: map[string]interface{}{"databaseURI": "postgres://username2:password2@localhost:5432/database_name2"},
		},
		{
			name: "Positive case: Set flag -sk and env JWT_SECRET_KEY",
			args: []string{"-sk=Secret3"},
			envs: []string{"JWT_SECRET_KEY=Secret4"},
			want: map[string]interface{}{"jwtSecretKey": "Secret4"},
		},
		{
			name: "Positive case: Set flag -sk only",
			args: []string{"-sk=Secret5"},
			envs: nil,
			want: map[string]interface{}{"jwtSecretKey": "Secret5"},
		},
		{
			name: "Positive case: Set env JWT_SECRET_KEY only",
			args: nil,
			envs: []string{"JWT_SECRET_KEY=Secret6"},
			want: map[string]interface{}{"jwtSecretKey": "Secret6"},
		},
		{
			name: "Positive case: Set flag -ll and env LOG_LEVEL",
			args: []string{"-ll=warn"},
			envs: []string{"LOG_LEVEL=error"},
			want: map[string]interface{}{"logLevel": "error"},
		},
		{
			name: "Positive case: Set flag -ll only",
			args: []string{"-ll=fatal"},
			envs: nil,
			want: map[string]interface{}{"logLevel": "fatal"},
		},
		{
			name: "Positive case: Set env LOG_LEVEL only",
			args: nil,
			envs: []string{"LOG_LEVEL=panic"},
			want: map[string]interface{}{"logLevel": "panic"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if len(tc.args) > 0 {
				os.Args = append(os.Args, tc.args...)
			}

			if len(tc.envs) > 0 {
				for _, v := range tc.envs {
					e := strings.Split(v, "=")

					suite.Require().GreaterOrEqual(len(e), 2)

					if len(e) > 2 {
						e[1] = strings.Join(e[1:], "=")
						e = e[:2]
					}

					err := os.Setenv(e[0], e[1])
					suite.Require().NoError(err)
				}
			}

			config := initConfig()

			configFields := reflect.ValueOf(config)

			for k, want := range tc.want {
				field := configFields.FieldByName(k)
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'. Expected: '%+v'. Actual: '%+v'.", k, reflect.ValueOf(want), reflect.ValueOf(field))
			}
		})
	}
}

func TestFlagsSuite(t *testing.T) {
	suite.Run(t, new(FlagsTestSuite))
}
