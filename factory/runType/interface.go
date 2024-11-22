package runType

// RunTypeComponentsCreator is the interface for creating run type components
type RunTypeComponentsCreator interface {
	Create() *runTypeComponents
	IsInterfaceNil() bool
}
