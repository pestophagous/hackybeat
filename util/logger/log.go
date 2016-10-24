package log

type LogAdapter struct {
	Err   func(format string, v ...interface{})
	Warn  func(format string, v ...interface{})
	Info  func(format string, v ...interface{})
	Debug func(format string, v ...interface{})
}

func noop(format string, v ...interface{}) {
}

func NewNoopLogAdapter() *LogAdapter {
	l := &LogAdapter{
		Err:   noop,
		Warn:  noop,
		Info:  noop,
		Debug: noop,
	}

	return l
}

func NewNoopLogWithNilCheck() *LogWithNilCheck {
	return &LogWithNilCheck{NewNoopLogAdapter()}
}

type LogWithNilCheck struct {
	L *LogAdapter
}

type IdentifiableForLog interface {
	InstanceIdForLogging() string
}

// Our internal LogAdapter field most likely consists of a bunch of closures that are
// capturing referencs to who-knows-what.  We might at times like to release these references.
func (this *LogWithNilCheck) ReleaseLog() {
	this.L = nil
}

func (this *LogWithNilCheck) Err(format string, v ...interface{}) {
	if this.L != nil {
		this.L.Err(format, v...)
	}
}

func (this *LogWithNilCheck) Warn(format string, v ...interface{}) {
	if this.L != nil {
		this.L.Warn(format, v...)
	}
}

func (this *LogWithNilCheck) Info(format string, v ...interface{}) {
	if this.L != nil {
		this.L.Info(format, v...)
	}
}

func (this *LogWithNilCheck) Debug(format string, v ...interface{}) {
	if this.L != nil {
		this.L.Debug(format, v...)
	}
}

// convenience function if you're already inside a block with a proven non-nil error:
func (this *LogWithNilCheck) LogFailureOf(what string, caller IdentifiableForLog, e error) {
	this.Err("%s failed on %v. %v", what, caller.InstanceIdForLogging(), e)
}

// convenience function when an error may or may not be nil, but you only want to log when it's non-nil:
func (this *LogWithNilCheck) LogPossibleFailureOf(what string, caller IdentifiableForLog, e error) {
	if e != nil {
		this.LogFailureOf(what, caller, e)
	}
}
