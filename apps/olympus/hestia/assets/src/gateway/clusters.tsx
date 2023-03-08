import {zeusApi} from './axios/axios';

class ClustersApiGateway {
    async getClusters(): Promise<any>  {
        const url = `/v1/infra/read/org/topologies`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
            headers: {
                'Authorization': `Bearer ${sessionID}`
            }}
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
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                }}
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
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                    'CloudCtxNsID': `${clusterID}`
                }}
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
