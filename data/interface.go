package data

func (h *HyperblockApiResponse) Hash() string {
	return h.Data.Hyperblock.Hash
}

func (h *HyperblockApiResponse) Nonce() uint64 {
	return h.Data.Hyperblock.Nonce
}

func (h *BlockApiResponse) Hash() string {
	return h.Data.Block.Hash
}

func (h *BlockApiResponse) Nonce() uint64 {
	return h.Data.Block.Nonce
}
