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
    async checkAuth(): Promise<any>  {
        const url = `/v1/auth/status`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await hestiaApi.get(url, config)
    }
    async startPlatformAuthFlow(platformName: string): Promise<any> {
        const url = `/social/v1/auth/${platformName}/callback`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await hestiaApi.get(url, config)
    }

    async callbackPlatformAuthFlow(platformName: string, code: string, state: string) {
        const url = `/social/v1/${platformName}/callback?code=${encodeURIComponent(code)}&state=${encodeURIComponent(state)}`;
        // const sessionID = inMemoryJWT.getToken();
        // let config = {
        //     headers: {
        //         'Authorization': `Bearer ${sessionID}`
        //     },
        //     withCredentials: true,
        // }
        return await hestiaApi.get(url);
    }
}
export const accessApiGateway = new AccessApiGateway();

