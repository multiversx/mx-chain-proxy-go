package runType

// RunTypeComponentsCreator is the interface for the runTypeComponentsCreator
type RunTypeComponentsCreator interface {
	Create() *runTypeComponents
	IsInterfaceNil() bool
}
