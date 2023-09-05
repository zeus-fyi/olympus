package v1_iris

import (
	"github.com/labstack/echo/v4"
	iris_catalog_procedures "github.com/zeus-fyi/olympus/iris/api/v1/procedures"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

func GetProcedureTemplateJsonRPC(rgName, procName string, req *iris_api_requests.ApiProxyRequest, stageTwoPayload echo.Map) (iris_programmable_proxy_v1_beta.IrisRoutingProcedure, error) {
	switch procName {
	case iris_catalog_procedures.EthMaxBlockAggReduce:
		return GetEthMaxBlockAggReduceTemplate(rgName, req, stageTwoPayload)
	case iris_catalog_procedures.AvaxMaxBlockAggReduce:
		return GetAvaxMaxBlockAggReduceTemplate(rgName, req, stageTwoPayload)
	case iris_catalog_procedures.NearMaxBlockAggReduce:
		return GetNearMaxBlockAggReduceTemplate(rgName, req, stageTwoPayload)
	case iris_catalog_procedures.BtcMaxBlockAggReduce:
		return GetBtcMaxBlockAggReduceTemplate(rgName, req, stageTwoPayload)
	default:
		return iris_programmable_proxy_v1_beta.IrisRoutingProcedure{}, nil
	}
}

func GetEthMaxBlockAggReduceTemplate(rgName string, req *iris_api_requests.ApiProxyRequest, stageTwoPayload echo.Map) (iris_programmable_proxy_v1_beta.IrisRoutingProcedure, error) {
	fnRule := iris_programmable_proxy_v1_beta.FanInRuleFirstValidResponse
	ph := ProcedureHeaders{
		XAggOp:               "max",
		XAggKey:              "result",
		XAggKeyValueDataType: "int",
		XAggFilterFanIn:      &fnRule,
		ForwardPayload:       stageTwoPayload,
	}
	req.Payload = iris_catalog_procedures.ProcedureStageOnePayload(iris_catalog_procedures.EthMaxBlockAggReduce)
	return ph.GetGeneratedProcedure(rgName, req)
}

func GetNearMaxBlockAggReduceTemplate(rgName string, req *iris_api_requests.ApiProxyRequest, stageTwoPayload echo.Map) (iris_programmable_proxy_v1_beta.IrisRoutingProcedure, error) {
	fnRule := iris_programmable_proxy_v1_beta.FanInRuleFirstValidResponse
	ph := ProcedureHeaders{
		XAggOp:               "max",
		XAggKey:              "result,sync_info,latest_block_height",
		XAggKeyValueDataType: "int",
		XAggFilterFanIn:      &fnRule,
		ForwardPayload:       stageTwoPayload,
	}
	req.Payload = iris_catalog_procedures.ProcedureStageOnePayload(iris_catalog_procedures.NearMaxBlockAggReduce)
	return ph.GetGeneratedProcedure(rgName, req)
}

func GetBtcMaxBlockAggReduceTemplate(rgName string, req *iris_api_requests.ApiProxyRequest, stageTwoPayload echo.Map) (iris_programmable_proxy_v1_beta.IrisRoutingProcedure, error) {
	fnRule := iris_programmable_proxy_v1_beta.FanInRuleFirstValidResponse
	ph := ProcedureHeaders{
		XAggOp:               "max",
		XAggKey:              "result,sync_info,latest_block_height",
		XAggKeyValueDataType: "int",
		XAggFilterFanIn:      &fnRule,
		ForwardPayload:       stageTwoPayload,
	}
	req.Payload = iris_catalog_procedures.ProcedureStageOnePayload(iris_catalog_procedures.AvaxMaxBlockAggReduce)
	return ph.GetGeneratedProcedure(rgName, req)
}

func GetAvaxMaxBlockAggReduceTemplate(rgName string, req *iris_api_requests.ApiProxyRequest, stageTwoPayload echo.Map) (iris_programmable_proxy_v1_beta.IrisRoutingProcedure, error) {
	fnRule := iris_programmable_proxy_v1_beta.FanInRuleFirstValidResponse
	ph := ProcedureHeaders{
		XAggOp:               "max",
		XAggKey:              "result,sync_info,latest_block_height",
		XAggKeyValueDataType: "int",
		XAggFilterFanIn:      &fnRule,
		ForwardPayload:       stageTwoPayload,
	}
	req.Payload = iris_catalog_procedures.ProcedureStageOnePayload(iris_catalog_procedures.AvaxMaxBlockAggReduce)
	return ph.GetGeneratedProcedure(rgName, req)
}
