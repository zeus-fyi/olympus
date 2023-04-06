import {zeusApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";
import {TopologySystemComponents} from "../redux/apps/apps.types";
import {ClusterPreview} from "../redux/clusters/clusters.types";

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
    async getPrivateAppDetails(id: string): Promise<ClusterPreview>  {
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
            return {} as ClusterPreview
        }
    }
}

export const appsApiGateway = new AppsApiGateway();
