package contracts

type Block = struct {
	Body              string
	Hash              string
	Height            int
	Owner             string
	PreviousBlockHash string
	Time              int64
}

type StarData struct {
	Address   string
	Message   string
	Data      []byte
	Signature string
}

type BlockchainOperator interface {
	RequestMessageOwnershipVerification(addr string) (string, error)
	GetBlockByHeight(h int) (Block, error)
	GetBlockByHash(h string) (Block, error)
	GetStarsByWalletAddress(addr string) []string
	SubmitStar(star StarData) (Block, error)
	Validate() (bool, []string)
}

type Clock interface {
	GetTime() int64
}
