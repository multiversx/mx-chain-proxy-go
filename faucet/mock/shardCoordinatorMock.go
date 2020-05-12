package mock

type ShardCoordinatorMock struct {
}

func (scm *ShardCoordinatorMock) NumberOfShards() uint32 {
	panic("implement me")
}

func (scm *ShardCoordinatorMock) ComputeId(_ []byte) uint32 {
	return uint32(1)
}

func (scm *ShardCoordinatorMock) SetSelfId(_ uint32) error {
	panic("implement me")
}

func (scm *ShardCoordinatorMock) SelfId() uint32 {
	return 0
}

func (scm *ShardCoordinatorMock) SameShard(_, _ []byte) bool {
	return true
}

func (scm *ShardCoordinatorMock) CommunicationIdentifier(_ uint32) string {
	return "0_1"
}

// IsInterfaceNil returns true if there is no value under the interface
func (scm *ShardCoordinatorMock) IsInterfaceNil() bool {
	return scm == nil
}
