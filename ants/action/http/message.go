package http

// welcome struct
type WelcomeInfo struct {
	Message  string
	Greeting string
	Time     string
}

// result of start spider
type StartSpiderResult struct {
	Success bool
	Message string
	Spider  string
	Time    string
}
