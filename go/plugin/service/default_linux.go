package service

func newDefaultBindOption() *bindOption {
	result := &bindOption{}
	WithAnonymousPipe()(result)
	return result
}
