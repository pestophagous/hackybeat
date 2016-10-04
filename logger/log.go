package log

type LogAdapter struct {
	Err   func(format string, v ...interface{})
	Warn  func(format string, v ...interface{})
	Info  func(format string, v ...interface{})
	Debug func(format string, v ...interface{})
}

type LogWithNilCheck struct {
	L *LogAdapter
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
