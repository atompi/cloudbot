package options

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Version string = "v0.0.1"

type LogOptions struct {
	Level    string `yaml:"level"`
	Path     string `yaml:"path"`
	MaxSize  int    `yaml:"maxsize"`
	MaxAge   int    `yaml:"maxage"`
	Compress bool   `yaml:"compress"`
}

type CoreOptions struct {
	Log LogOptions `yaml:"log"`
}

type AliyunOptions struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	RegionId        string `yaml:"region_id"`
	Endpoint        string `yaml:"endpoint"`
}

type TencentOptions struct {
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
	Region    string `yaml:"region"`
	Endpoint  string `yaml:"endpoint"`
}

type InputOutputOptions struct {
	Type   string `yaml:"type"`
	Path   string `yaml:"path"`
	Target string `yaml:"target"`
}

type TaskOptions struct {
	Name    string             `yaml:"name"`
	Enabled bool               `yaml:"enabled"`
	Type    string             `yaml:"type"`
	Threads int                `yaml:"threads"`
	Aliyun  AliyunOptions      `yaml:"aliyun"`
	Tencent TencentOptions     `yaml:"tencent"`
	Input   InputOutputOptions `yaml:"input"`
	Output  InputOutputOptions `yaml:"output"`
}

type Options struct {
	Core  CoreOptions   `yaml:"core"`
	Tasks []TaskOptions `yaml:"tasks"`
}

func NewOptions() (opts Options) {
	optsSource := viper.AllSettings()
	err := createOptions(optsSource, &opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "create options failed:", err)
		os.Exit(1)
	}
	return
}
