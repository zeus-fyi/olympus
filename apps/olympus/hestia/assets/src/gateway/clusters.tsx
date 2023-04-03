import {zeusApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";

class ClustersApiGateway {
    async previewCreateCluster(params: any): Promise<any>  {
        const url = `/v1/infra/ui/preview/create`;
        console.log(params)

        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload = {
                'cluster': params
            }
            return await zeusApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending cluster preview create request');
            console.error(exc);
            return
        }
    }
    async createCluster(params: any): Promise<any>  {
        const url = `/v1/infra/ui/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload = {
                'cluster': params
            }
            return await zeusApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending cluster create request');
            console.error(exc);
            return
        }
    }
    async getClusters(): Promise<any>  {
        const url = `/v1/infra/read/org/topologies`;
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
            console.error('error sending cluster get request');
            console.error(exc);
            return
        }
    }
    async getClusterTopologies(params: any): Promise<any>  {
        const url = `/v1/deploy/cluster/status`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload = {
                'cloudCtxNsID': `${params.id}`
            }
            return await zeusApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending get cluster topologies at cloud ctx ns request');
            console.error(exc);
            return
        }
    }
    async getClusterPodsAudit(clusterID: any): Promise<any>  {
        const url = `/v1/pods`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                    'CloudCtxNsID': `${clusterID}`
                },
                withCredentials: true,
            }
            const payload = {
                action: "describe-audit",
            }
            return await zeusApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending get cluster topologies at cloud ctx ns request');
            console.error(exc);
            return
        }
    }
}
export const clustersApiGateway = new ClustersApiGateway();
