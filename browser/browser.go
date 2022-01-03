package browser

type Launcher interface {
	Open(url string) error
}

type launcher struct{}

func NewLauncher() Launcher {
	return launcher{}
}

func (_ launcher) Open(url string) error {
	return getBrowserCmd(url).Start()
}
