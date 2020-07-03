package server

//Server ...
type Server interface {
	Serve(addr string)
	Stop()
}
