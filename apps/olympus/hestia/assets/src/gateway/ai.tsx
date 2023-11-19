import {zeusApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";

class AiApiGateway {
    async searchRequest(params: any): Promise<any> {
        const url = `/v1/search`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload = {
                'searchParams': params
            }
            return await zeusApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending search request');
            console.error(exc);
            return
        }
    }
}

export const aiApiGateway = new AiApiGateway();
