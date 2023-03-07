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
           console.log(config);
            console.log(zeusApi)
            return await zeusApi.get(url, config)
        } catch (exc) {
            console.error('error sending cluster get request');
            console.error(exc);
            return
        }
    }
}
export const clustersApiGateway = new ClustersApiGateway();

