package config

import (
	"flag"
	"github.com/stretchr/testify/suite"
	"os"
	"reflect"
	"strings"
	"testing"
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
				"runAddr":     "localhost:8080",
				"accrualAddr": "localhost:8081",
				"databaseURI": "",
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
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'", k)
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
				"runAddr":     "",
				"accrualAddr": "",
				"databaseURI": "",
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
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'", k)
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
				"runAddr":     "localhost:8080",
				"accrualAddr": "localhost:8081",
				"databaseURI": "",
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
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'", k)
			}
		})
	}
}

func TestFlagsSuite(t *testing.T) {
	suite.Run(t, new(FlagsTestSuite))
}
