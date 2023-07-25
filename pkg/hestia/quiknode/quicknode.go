package quicknode

/*

	If referers is set, you will need to pass it as a header with each request you make to the customer's endpoint.
	Additionally, you should definitely store all of these pieces of information.

	We will use the access-url to show the customer how they can access the service
	(for example, if it’s an explorer, graph instance, some kind of index, GraphQL API or, REST API).
	If you are providing a JSON-RPC method, then simply return access-url as null. We will use the dashboard-url
	to log the customer into your service with SSO. We have a guide on SSO [here]().

	If responding with error, please be sure to use a non-200 HTTP status code.
*/
