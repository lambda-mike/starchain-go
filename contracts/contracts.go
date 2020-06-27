package contracts

type Blockchain interface {
	RequestMessageOwnershipVerification(addr string) (string, error)
}
