import {irisApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";

class IrisLoadBalancingApiGateway {
    async sendJsonRpcRequest(routeGroup: string, payload: any): Promise<any>  {
        const url = `/v1/router`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                    'X-Route-Group': `${routeGroup}`
                },
                withCredentials: true,
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

