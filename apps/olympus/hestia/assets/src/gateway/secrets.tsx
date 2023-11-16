import {hestiaApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";

const config = {
    withCredentials: true,
};

export interface SecretsRequest {
    name: string;
    key: string;
    value: string;
}

class SecretsApiGateway {
    async upsertSecret(request: SecretsRequest): Promise<any> {
        const url = `/v1/secrets/upsert`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.post(url, request, config);
        } catch (exc) {
            console.error('error upserting secret');
            console.error(exc);
        }
    }

    async deleteSecret(request: SecretsRequest): Promise<any> {
        const url = `/v1/secret/${request.name}`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.delete(url, config);
        } catch (exc) {
            console.error('error upserting secret');
            console.error(exc);
        }
    }

    async getSecrets(): Promise<any> {
        const url = `/v1/secrets`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.get(url, config);
        } catch (exc) {
            console.error('error getting secrets');
            console.error(exc);
        }
    }

    async getSecret(ref: string): Promise<any> {
        const url = `/v1/secret/${ref}`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.get(url, config);
        } catch (exc) {
            console.error(`error getting secret: ${ref}`);
            console.error(exc);
        }
    }
}

export const secretsApiGateway = new SecretsApiGateway();