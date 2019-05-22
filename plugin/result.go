package plugin

type Result struct {
	Safety   bool
	Target   Target
	Plugin   *Plugin
	Request  string
	Response string
}

func NewResult(safety bool, target Target, plugin *Plugin) Result {
	return Result{
		Safety: safety,
		Target: target,
		Plugin: plugin,
	}
}

func NewSafety(target Target, plugin *Plugin) Result {
	return NewResult(true, target, plugin)
}

func NewUnSafety(target Target, plugin *Plugin) Result {
	return NewResult(false, target, plugin)
}
