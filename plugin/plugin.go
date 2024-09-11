package plugin

type SpecData any

type Plugin interface {
	LoadSpecFile(path string) SpecData
}
