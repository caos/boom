package name

type Application string

func (a Application) String() string {
	return string(a)
}

type Templator string

func (t Templator) String() string {
	return string(t)
}
