package options

type CloudProviderOptions struct {
	AccessKeyId     string
	AccessKeySecret string
	RegionId        string
	Endpoint        string
}

type InputOutputOptions struct {
	Type   string
	Path   string
	Target string
}

type TaskOptions struct {
	Name          string
	Enabled       bool
	Type          string
	Threads       int
	CloudProvider CloudProviderOptions
	Input         InputOutputOptions
	Output        InputOutputOptions
}
