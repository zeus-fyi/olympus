import {hestiaApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";

class LoadBalancingApiGateway {
    async getEndpoints(): Promise<any>  {
        const url = `/iris/routes/read`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.get(url, config)
        } catch (exc) {
            console.error('error sending get customer endpoints request');
            console.error(exc);
            return
        }
    }
    async createEndpoints(): Promise<any>  {
        const url = `/iris/routes/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.post(url, config)
        } catch (exc) {
            console.error('error sending get customer endpoints request');
            console.error(exc);
            return
        }
    }
    async deleteEndpoints(): Promise<any>  {
        const url = `/iris/routes/delete`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.post(url, config)
        } catch (exc) {
            console.error('error sending get customer endpoints request');
            console.error(exc);
            return
        }
    }
    async deleteRoutingGroupEndpoints(groupName: string): Promise<any>  {
        const url = `/iris/routes/delete`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.post(url, config)
        } catch (exc) {
            console.error('error sending get customer endpoints request');
            console.error(exc);
            return
        }
    }
}
export const loadBalancingApiGateway = new LoadBalancingApiGateway();
