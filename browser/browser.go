package browser

// Launcher - interface for opening a browser window
type Launcher interface {
	Open(url string) error
}

type launcher struct{}

// NewLauncher instantiates a struct implementing the Launcher interface
func NewLauncher() Launcher {
	return launcher{}
}

// Open will launch the default browser with the given URL
func (launcher) Open(url string) error {
	return getBrowserCmd(url).Start()
}
