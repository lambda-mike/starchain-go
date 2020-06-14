package block

type Block struct {
	data string
}

func NewBlock(data string) *Block {
	return &Block{data}
}
