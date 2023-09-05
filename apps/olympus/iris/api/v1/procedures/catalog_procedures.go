package iris_catalog_procedures

const (
	eth_maxBlockAggReduce  = "eth_maxBlockAggReduce"
	avax_maxBlockAggReduce = "avax_maxBlockAggReduce"
	near_maxBlockAggReduce = "near_maxBlockAggReduce"
	btc_maxBlockAggReduce  = "btc_maxBlockAggReduce"
)

func Procedure(procName string) {
	switch procName {
	case eth_maxBlockAggReduce:
		// do something
	case avax_maxBlockAggReduce:
		// do something
	case near_maxBlockAggReduce:
		// do something
	case btc_maxBlockAggReduce:
		// do something
	default:
		// do something
	}
	return
}
