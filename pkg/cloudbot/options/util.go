package options

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func setOptions(optsSource, opts any) error {
	optsYaml, err := yaml.Marshal(optsSource)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(optsYaml, opts)
	return err
}

func New() (opts Options) {
	optsSource := viper.AllSettings()
	err := setOptions(optsSource, &opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "create options failed:", err)
		os.Exit(1)
	}
	return
}
