package functions

type Function struct {
	Id, Image string

	Environment map[string]string
}

type FunctionBuild struct {
	Image       string
	Environment map[string]string
}
