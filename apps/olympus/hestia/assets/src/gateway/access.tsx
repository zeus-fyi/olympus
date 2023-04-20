import {hestiaApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";

class AccessApiGateway {
    async sendApiKeyGenRequest(): Promise<any>  {
        const url = `/v1/api/key/create`;
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
            console.error('error sending create api key request');
            console.error(exc);
            return
        }
    }
}
export const accessApiGateway = new AccessApiGateway();

