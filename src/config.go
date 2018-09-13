
package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type config struct {
	TrySystemResolver		bool		`toml:"try_system_resolver"`
	CustomResolvers			[]string		`toml:"custom_resolvers"`
	Record					string		`toml:"record"`
	RecordType				string		`toml:"record_type"`
	ExpectedValue			string		`toml:"expected_value"`
	TriggerValue			string		`toml:"trigger_value"`
	DeleteFiles				[]string	`toml:"delete_files"`
	ExecuteScripts			[]string	`toml:"execute_scripts"`
	TriggerOnUncertain		bool		`toml:"trigger_on_uncertain"`
	MaxUncertainTolerance	uint		`toml:"max_uncertain_tolerance"`
	CheckInterval			uint		`toml:"check_interval"`
	ExitAfterTrigger		bool		`toml:"exit_after_trigger"`
}

func loadConfig(path string) (*config, error) {
	conf := &config{}
	metaData, err := toml.DecodeFile(path, conf)
	if err != nil {
		return nil, err
	}
	for _, key := range metaData.Undecoded() {
		return nil, &configError{fmt.Sprintf("unknown option %q", key.String())}
	}

	if conf.CheckInterval == 0 {
		conf.CheckInterval = 60
	}

	return conf, nil
}

type configError struct {
	err string
}

func (e *configError) Error() string {
	return e.err
}