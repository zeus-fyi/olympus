import {zeusApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";
import {PostWorkflowsActionRequest, PostWorkflowsRequest, Retrieval, TaskModelInstructions} from "../redux/ai/ai.types";

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
    async analyzeSearchRequest(params: any): Promise<any> {
        const url = `/v1/search/analyze`;
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
    async createAiWorkflowRequest(params: PostWorkflowsRequest): Promise<any> {
        const url = `/v1/workflows/ai`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await zeusApi.post(url, params, config)
        } catch (exc) {
            console.error('error sending search request');
            console.error(exc);
            return
        }
    }
    async createOrUpdateTaskRequest(params: TaskModelInstructions): Promise<any> {
        const url = `/v1/tasks/ai`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await zeusApi.post(url, params, config)
        } catch (exc) {
            console.error('error sending task create or update request');
            console.error(exc);
            return
        }
    }
    async createOrUpdateRetrieval(params: Retrieval): Promise<any> {
        const url = `/v1/retrievals/ai`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await zeusApi.post(url, params, config)
        } catch (exc) {
            console.error('error sending task create or update request');
            console.error(exc);
            return
        }
    }
    async getWorkflowsRequest(): Promise<any> {
        const url = `/v1/workflows/ai`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await zeusApi.get(url, config)
        } catch (exc) {
            console.error('error sending search request');
            console.error(exc);
            return
        }
    }
    async execWorkflowsActionRequest(params: PostWorkflowsActionRequest): Promise<any> {
        const url = `/v1/workflows/ai/start`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await zeusApi.post(url, params, config)
        } catch (exc) {
            console.error('error sending task create or update request');
            console.error(exc);
            return
        }
    }


}

export const aiApiGateway = new AiApiGateway();
