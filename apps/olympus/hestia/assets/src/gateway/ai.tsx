import {zeusApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";
import {
    AiSearchParams,
    DeleteWorkflowsActionRequest,
    PostCreateOrUpdateSearchIndexerRequest,
    PostRunsActionRequest,
    PostSearchIndexerActionsRequest,
    PostWorkflowsActionRequest,
    PostWorkflowsRequest,
    TaskModelInstructions,
} from "../redux/ai/ai.types";
import {Assistant, Retrieval} from "../redux/ai/ai.types.retrievals";
import {JsonSchemaDefinition} from "../redux/ai/ai.types.schemas";
import {EvalFn} from "../redux/ai/ai.types.evals";
import {TriggerAction, TriggerActionApprovalPutRequest} from "../redux/ai/ai.types.triggers";

class AiApiGateway {
    async searchRequest(params: AiSearchParams): Promise<any> {
        const url = `/v1/search`;
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
    }
    async createAiWorkflowRequest(params: PostWorkflowsRequest): Promise<any> {
        const url = `/v1/workflows/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
    async createOrUpdateJsonSchema(params: JsonSchemaDefinition): Promise<any> {
        const url = `/v1/schemas/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
    async createOrUpdateTaskRequest(params: TaskModelInstructions): Promise<any> {
        const url = `/v1/tasks/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
    async createOrUpdateRetrieval(params: Retrieval): Promise<any> {
        const url = `/v1/retrievals/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
    async createOrUpdateAssistant(params: Assistant): Promise<any> {
        const url = `/v1/assistants/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
    async createOrUpdateAction(params: TriggerAction): Promise<any> {
        const url = `/v1/actions/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
    async getActions(): Promise<any> {
        const url = `/v1/actions/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.get(url, config)
    }
    async getRuns(): Promise<any> {
        const url = `/v1/runs/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.get(url, config)
    }
    async getRun(rid: string): Promise<any> {
        const url = `/v1/run/ai/${rid}`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.get(url, config)
    }
    async getRunsUI(): Promise<any> {
        const url = `/v1/runs/ai/ui`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.get(url, config)
    }
    async getTasks(): Promise<any> {
        const url = `/v1/tasks/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.get(url, config)
    }
    async updateActionApproval(params: TriggerActionApprovalPutRequest): Promise<any> {
        const url = `/v1/actions/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.put(url, params, config)
    }
    async createOrUpdateEval(params: EvalFn): Promise<any> {
        const url = `/v1/evals/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
    async getWorkflowsRequest(): Promise<any> {
        const url = `/v1/workflows/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.get(url, config)
    }
    async deleteWorkflowsActionRequest(params: DeleteWorkflowsActionRequest): Promise<any> {
        const url = `/v1/workflows/ai`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
            data: params
        }
        return await zeusApi.delete(url, config)
    }
    async execWorkflowsActionRequest(params: PostWorkflowsActionRequest): Promise<any> {
        const url = `/v1/workflows/ai/actions`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
    async execRunsActionRequest(params: PostRunsActionRequest): Promise<any> {
        const url = `/v1/runs/ai/actions`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
    async searchIndexerCreateOrUpdateRequest(params: PostCreateOrUpdateSearchIndexerRequest): Promise<any> {
        const url = `/v1/search/indexer`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
    async searchIndexerCreateOrUpdateActionRequest(params: PostSearchIndexerActionsRequest): Promise<any> {
        const url = `/v1/search/indexer/actions`;
        const sessionID = inMemoryJWT.getToken();
        let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            },
            withCredentials: true,
        }
        return await zeusApi.post(url, params, config)
    }
}

export const aiApiGateway = new AiApiGateway();
