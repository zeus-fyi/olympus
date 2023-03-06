import axios from './axios/axios';

const config = {
};

class ClustersApiGateway {
    async getClusters(user: string): Promise<any>  {
        const url = `clusters`;
        try {
            return await axios.post(url, {
                user: user,
            }, config)
        } catch (exc) {
            console.error('error sending clusters fetch request');
            console.error(exc);
            return
        }
    }
}
export const clustersApiGateway = new ClustersApiGateway();

