package forecastle

// App struct that contains information about an app that is exposed to forecastle
type App struct {
	Name     string
	Icon     string
	Group    string
	URL      string
	IsCustom bool // To be used later by the UI to indicate whether the App was generated or manually added
}
