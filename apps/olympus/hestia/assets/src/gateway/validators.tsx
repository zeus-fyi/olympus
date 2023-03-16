import {hestiaApi} from './axios/axios';
import {AwsCredentialIdentity} from "@aws-sdk/types/dist-types/identity";

class ValidatorsApiGateway {
    async getValidators(): Promise<any>  {
        const url = `/v1/ethereum/validators/service/info`;
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
    async verifyValidators(credentials: AwsCredentialIdentity, fnUrl: string, depositSlice: [{}]): Promise<any>  {
        const url = `/v1/ethereum/validators/service/create`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            console.log(fnUrl)
            console.log(depositSlice)
            console.log(credentials)
            return await hestiaApi.post(url, depositSlice, config)
        } catch (exc) {
            console.error('error sending create lambda function keystores layer');
            console.error(exc);
            return
        }
    }
    async createValidatorsServiceRequest(payload: any): Promise<any>  {
        const url = `/v1/ethereum/validators/service/create`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create lambda function keystores layer');
            console.error(exc);
            return
        }
    }
}
export const validatorsApiGateway = new ValidatorsApiGateway();

