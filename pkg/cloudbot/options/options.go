package options

var Version string = "v0.0.1"

type LogOptions struct {
	Level      string `yaml:"level"`
	Path       string `yaml:"path"`
	Format     string `yaml:"format"`
	MaxSize    int    `yaml:"maxsize"`
	MaxAge     int    `yaml:"maxage"`
	MaxBackups int    `yaml:"maxbackups"`
	Compress   bool   `yaml:"compress"`
	MultiFiles bool   `yaml:"multifiles"`
}

type CoreOptions struct {
	Log            LogOptions `yaml:"log"`
	TasksXlsxFile  string     `yaml:"tasksxlsxfile"`
	TasksXlsxSheet string     `yaml:"tasksxlsxsheet"`
}

type Options struct {
	Core CoreOptions `yaml:"core"`
}
