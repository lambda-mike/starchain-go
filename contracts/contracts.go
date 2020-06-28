package contracts

type Block = struct {
	Body              string
	Hash              string
	Height            int
	Owner             string
	PreviousBlockHash string
	Time              int64
}

type Blockchain interface {
	RequestMessageOwnershipVerification(addr string) (string, error)
	GetBlockByHeight(h int) (Block, error)
}
