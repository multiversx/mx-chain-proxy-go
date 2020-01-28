package ring

// ObserversRingHandler defines what a ring of observers should be able to do
type ObserversRingHandler interface {
	Next() string
	Len() int
	IsInterfaceNil() bool
}
