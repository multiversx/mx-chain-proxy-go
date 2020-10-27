package mock

type ShardCoordinatorMock struct {
	NumShards uint32
}

func (scm *ShardCoordinatorMock) NumberOfShards() uint32 {
	return scm.NumShards
}

func (scm *ShardCoordinatorMock) ComputeId(_ []byte) uint32 {
	return uint32(1)
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
