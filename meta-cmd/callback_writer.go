package metacmd

type CallbackWriter struct {
	Callback func(p []byte) (n int, err error)
}

func (lw *CallbackWriter) Write(p []byte) (n int, err error) {
	if lw.Callback != nil {
		return lw.Callback(p)
	}
	return len(p), nil
}
