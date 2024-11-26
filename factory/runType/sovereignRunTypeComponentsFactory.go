package runType

type sovereignRunTypeComponentsFactory struct{}

// NewSovereignRunTypeComponentsFactory will return a new instance of sovereign run type components factory
func NewSovereignRunTypeComponentsFactory() *sovereignRunTypeComponentsFactory {
	return &sovereignRunTypeComponentsFactory{}
}

// Create will create the run type components
func (srtcf *sovereignRunTypeComponentsFactory) Create() *runTypeComponents {
	return &runTypeComponents{}
}

// IsInterfaceNil returns true if there is no value under the interface
func (srtcf *sovereignRunTypeComponentsFactory) IsInterfaceNil() bool {
	return srtcf == nil
}
