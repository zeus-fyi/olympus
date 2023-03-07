import {hestiaApi} from './axios/axios';

class ValidatorsApiGateway {
    async getValidators(network: string): Promise<any>  {
        const url = `/v1/validators/service/info/`+network;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                }}
            return await hestiaApi.get(url, config)
        } catch (exc) {
            console.error('error sending get validators request');
            console.error(exc);
            return
        }
    }
}
export const validatorsApiGateway = new ValidatorsApiGateway();

