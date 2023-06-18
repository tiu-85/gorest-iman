package values

type Args struct {
	ConfigFilename string `short:"c" long:"config" description:"A path to config file" required:"true"`
	RunDry         bool   `short:"r" long:"run-dry" description:"Run-Dry mod. If exit code 0 deps graph is valid" required:"false"`
}
