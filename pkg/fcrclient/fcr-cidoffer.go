package fcrclient


// PieceCIDOffer contains offer information 
type PieceCIDOffer struct {
	pieceCID []byte
	providerID []byte
	price    int64
	expiry   int64
	expectedLatency int64
}