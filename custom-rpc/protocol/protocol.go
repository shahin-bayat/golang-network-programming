package protocol

type Request struct {
	ServiceMethod string
	Args          any
}

type Response struct {
	Reply any
	Error string
}
