import {irisApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";

class IrisLoadBalancingApiGateway {
    async sendJsonRpcRequest(routeGroup: string, payload: any, planName: string): Promise<any>  {
        const url = `/v1/router`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                    'X-Route-Group': `${routeGroup}`,
                    // 'X-Load-Balancing-Strategy': 'Adaptive',
                    // 'X-Adaptive-Metrics-Key': 'JSON-RPC',
                    'Content-Type': 'application/json'
                },
                withCredentials: true,
            }
            if (planName === "lite") {
                let config = {
                    headers: {
                        'Authorization': `Bearer ${sessionID}`,
                        'X-Route-Group': `${routeGroup}`,
                        'X-Load-Balancing-Strategy': 'RoundRobin',
                        'Content-Type': 'application/json'
                    },
                    withCredentials: true,
                }
                return await irisApi.post(url, payload, config)
            }
            return await irisApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending endpoints table request');
            console.error(exc);
            return
        }
    }
}
export const IrisApiGateway = new IrisLoadBalancingApiGateway();

