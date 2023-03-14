import {awsLambdaInvoke} from './axios/axios';
import {AwsRequest} from "./aws";

class AwsLambdaApiGateway {
    async invokeValidatorSecretsGeneration(url:string, ak: string, sk: string, mnemonicHdPwSecretName: string, ageKeySecretName: string): Promise<any> {
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }
            }
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: ak,
                    secretKey: sk,
                },
            };
            return await awsLambdaInvoke.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create new internal lambda user');
            console.error(exc);
            return
        }
    }
}

export const awsLambdaApiGateway = new AwsLambdaApiGateway();
