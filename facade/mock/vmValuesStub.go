package mock

type VmValuesProcessorStub struct {
	GetVmValueCalled func(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error)
}

func (gvps *VmValuesProcessorStub) GetVmValue(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
	return gvps.GetVmValueCalled(resType, address, funcName, argsBuff...)
}
