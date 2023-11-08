import {hestiaApi, zeusApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";
import {Cluster} from "../redux/clusters/clusters.types";

class ResourcesApiGateway {
    async getResources(): Promise<any>  {
        const url = `/v1/resources`;
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
            console.error('error sending get customer resources request');
            console.error(exc);
            return
        }
    }
    async searchNodeResources(ns: any): Promise<any>  {
        const url = `/v1/search/nodes`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await zeusApi.post(url, config)
        } catch (exc) {
            console.error('error sending get customer resources request');
            console.error(exc);
            return
        }
    }
    async destroyAppResource(orgResourceID: number): Promise<any>  {
        const url = `/v1/resources/destroy`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload = {
                orgResourceID: orgResourceID,
            }
            return await zeusApi.post(url, payload,config)
        } catch (exc) {
            console.error('error sending delete customer node resources request');
            console.error(exc);
            return []
        }
    }
    async destroyDeploy(cloudCtxNs: CloudCtxNs): Promise<any>  {
        const url = `/v1/deploy/ui/destroy`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload = {
                cloudProvider: cloudCtxNs.cloudProvider,
                region: cloudCtxNs.region,
                context: cloudCtxNs.context,
                namespace: cloudCtxNs.namespace,
            }
            return await zeusApi.post(url, payload,config)
        } catch (exc) {
            console.error('error sending delete customer node resources request');
            console.error(exc);
            return []
        }
    }
    async getAppResources(cluster: Cluster): Promise<any>  {
        const url = `/v1/nodes`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload: ActionRequest = {
                action: 'list',
                labels: {
                    'app': cluster.clusterName
                },
            }
            return await zeusApi.post(url, payload,config)
        } catch (exc) {
            console.error('error sending get customer app resources request');
            console.error(exc);
            return []
        }
    }
}
export const resourcesApiGateway = new ResourcesApiGateway();

interface ActionRequest {
    labels?: Record<string, string>;
    action: string;
}

export type CloudCtxNs = {
    cloudProvider: string;
    region: string;
    context: string;
    namespace: string;
};

export interface TopologyKubeCtxNs {
    topologyID: number;
    cloudCtxNs: CloudCtxNs;
}