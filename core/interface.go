package core

type RollupInter interface {
	RollupWithType(data []byte, daType int) ([]interface{}, error)
}

//type DAInter interface {
//	Store
//}
