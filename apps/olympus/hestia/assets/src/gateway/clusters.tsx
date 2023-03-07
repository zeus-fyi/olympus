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

    async getClusterTopologies(cloudCtxNsID: number): Promise<any>  {
        const url = `/v1/deploy/cluster/status`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                }}
            const payload = {
                cloudCtxNsID: cloudCtxNsID,
            }
            console.log(payload)
            return await zeusApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending get cluster topologies at cloud ctx ns request');
            console.error(exc);
            return
        }
    }
}
export const clustersApiGateway = new ClustersApiGateway();

