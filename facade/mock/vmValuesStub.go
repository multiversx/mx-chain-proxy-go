package mock

type VmValuesProcessorStub struct {
	GetVmValueCalled func(address string, funcName string, argsBuff ...[]byte) ([]byte, error)
}

func (gvps *VmValuesProcessorStub) GetVmValue(address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
	return gvps.GetVmValueCalled(address, funcName, argsBuff...)
}
