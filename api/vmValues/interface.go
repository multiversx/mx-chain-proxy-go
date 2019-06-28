package vmValues

// FacadeHandler interface defines methods that can be used from `elrondFacade` context variable
type FacadeHandler interface {
	GetVmValue(address string, funcName string, argsBuff ...[]byte) ([]byte, error)
}
