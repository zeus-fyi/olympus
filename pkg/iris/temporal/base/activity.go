package temporal_base

import "go.temporal.io/sdk/workflow"

/*
An Activity Definition is an exportable function or a struct method.
An Activity struct can have more than one method, with each method acting as a separate Activity Type.
Activities written as struct methods can use shared struct variables, such as:

	* an application level DB pool
	* client connection to another service
	* reusable utilities
	* any other expensive resources that you only want to initialize once per process
*/

type Activity struct {
	Name   string
	Params ActivityParams
}

type ActivityParams interface{}
type ActivityResponse interface{}
type ActivityDefinitionFn func(ctx workflow.Context, params ActivityParams) (ActivityResponse, error)
