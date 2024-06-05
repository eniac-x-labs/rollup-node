package core

type RollupInter interface {
	RollupWithType(data []byte, daType int) ([]interface{}, error)
	RetrieveFromDAWithType(daType int, args interface{}) ([]byte, error)
}

//type DAInter interface {
//	Store
//}
