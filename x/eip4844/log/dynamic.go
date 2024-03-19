package log

type LvlSetter interface {
	SetLogLevel(lvl Lvl)
}

// DynamicLogHandler allow runtime-configuration of the log handler.
type DynamicLogHandler struct {
	Handler // embedded, to expose any extra methods the underlying handler might provide
	maxLvl  Lvl
}

func NewDynamicLogHandler(lvl Lvl, h Handler) *DynamicLogHandler {
	return &DynamicLogHandler{
		Handler: h,
		maxLvl:  lvl,
	}
}

func (d *DynamicLogHandler) SetLogLevel(lvl Lvl) {
	d.maxLvl = lvl
}

func (d *DynamicLogHandler) Log(r *Record) error {
	if r.Lvl > d.maxLvl { // lower log level values are more critical
		return nil
	}
	return d.Handler.Log(r) // process the log
}
