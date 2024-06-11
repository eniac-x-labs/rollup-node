package nearda

type NearDAClient struct {
	//*near.Config
}

type INearDA interface {
	Store(data []byte) ([]byte, error)
	GetFromDA(frameRefBytes []byte, txIndex uint32) ([]byte, error)
}

//func NewNearDAClient(nearconf *NearDAConfig) (INearDA, error) {
//	conf, err := near.NewConfig(nearconf.Account, nearconf.Contract, nearconf.Key, nearconf.Network, nearconf.Ns)
//	if err != nil {
//		log.Error("NewConfig failed:", err)
//		return nil, err
//	}
//	return &NearDAClient{conf}, nil
//}
//func (n *NearDAClient) Store(data []byte) ([]byte, error) {
//	return n.ForceSubmit(data)
//}
//func (n *NearDAClient) GetFromDA(frameRefBytes []byte, txIndex uint32) ([]byte, error) {
//	return n.Get(frameRefBytes, txIndex)
//}

func NewNearDAClient(nearconf *NearDAConfig) (INearDA, error) {
	return nil, nil
}

func (n *NearDAClient) Store(data []byte) ([]byte, error) {
	return nil, nil
}
func (n *NearDAClient) GetFromDA(frameRefBytes []byte, txIndex uint32) ([]byte, error) {
	return nil, nil
}
