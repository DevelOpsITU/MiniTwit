package backgroundservices

var p program

func Init() {
	p = program{make(chan struct{})}
	start(&p)
}

func Dispose() {
	stop(&p)
}
