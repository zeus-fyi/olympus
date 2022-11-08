package deploy

//
//func HandleDeploymentActionRequest(c echo.Context) error {
//	request := new(DeploymentActionRequest)
//	if err := c.Bind(request); err != nil {
//		return err
//	}
//	if request.Action == "create" {
//		return request.DeployTopology(c)
//	}
//	if request.Action == "delete" {
//		return request.DeleteDeployedTopology(c)
//	}
//	return c.JSON(http.StatusBadRequest, nil)
//}
