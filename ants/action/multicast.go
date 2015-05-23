package action

type MulticastClient interface {
	send()
}

type MulticastServer interface {
	listen()
}
