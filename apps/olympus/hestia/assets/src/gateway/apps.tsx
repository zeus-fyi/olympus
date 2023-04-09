import {zeusApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";
import {Nodes, TopologySystemComponents} from "../redux/apps/apps.types";
import {Cluster, ClusterPreview} from "../redux/clusters/clusters.types";

class AppsApiGateway {
    async getPrivateApps(): Promise<TopologySystemComponents[]>  {
        const url = `/v1/infra/ui/private/apps`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await zeusApi.get(url, config).then((response) => {
                return response.data;
            })
        } catch (exc) {
            console.error('error sending get private apps request');
            console.error(exc);
            return []
        }
    }
    async getPrivateAppDetails(id: string): Promise<AppPageResponse>  {
        const url = `/v1/infra/ui/private/app/${id}`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await zeusApi.get(url, config).then((response) => {
                return response.data;
            })
        } catch (exc) {
            console.error('error sending get private apps request');
            console.error(exc);
            return {} as AppPageResponse
        }
    }
    async deployApp(payload: any): Promise<any>  {
        const url = `/v1/deploy/ui/app`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await zeusApi.post(url, payload, config).then((response) => {
                return response.data;
            })
        } catch (exc) {
            console.error('error sending app deployment request');
            console.error(exc);
            return {}
        }
    }
}

export const appsApiGateway = new AppsApiGateway();

export interface AppPageResponse {
    cluster: Cluster;
    clusterPreview: ClusterPreview;
    selectedComponentBaseName: string;
    selectedSkeletonBaseName: string;
    nodes: Nodes[];
}