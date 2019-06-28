package mock

type GetValuesProcessorStub struct {
	GetDataValueCalled func(address string, funcName string, argsBuff ...[]byte) ([]byte, error)
}

func (gvps *GetValuesProcessorStub) GetDataValue(address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
	return gvps.GetDataValueCalled(address, funcName, argsBuff...)
}
