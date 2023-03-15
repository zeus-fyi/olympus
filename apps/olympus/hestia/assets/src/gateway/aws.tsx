import {hestiaApi} from './axios/axios';

class AwsApiGateway {
    async createInternalLambdaUser(ak: string, sk: string): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/user/internal/lambda/create`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: ak,
                    secretKey: sk,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create new internal lambda user');
            console.error(exc);
            return
        }
    }
    async createExternalLambdaUser(ak: string, sk: string): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/user/external/lambda/create`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: ak,
                    secretKey: sk,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create new external lambda user');
            console.error(exc);
            return
        }
    }
    async createLambdaFunction(ak: string, sk: string): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/signer/create`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: ak,
                    secretKey: sk,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create lambda function');
            console.error(exc);
            return
        }
    }
    async createLambdaFunctionKeystoresLayer(ak: string, sk: string): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/keystore/create`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: ak,
                    secretKey: sk,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create lambda function keystores layer');
            console.error(exc);
            return
        }
    }
    async createValidatorsDepositDataLambda(ak: string, sk: string): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/deposits/create`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: ak,
                    secretKey: sk,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create validator deposits lambda function');
            console.error(exc);
            return
        }
    }
    async createValidatorsAgeEncryptedKeystoresZipLambda(ak: string, sk: string): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/keystores/zip/create`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: ak,
                    secretKey: sk,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create validator encrypted keystores lambda function');
            console.error(exc);
            return
        }
    }
    async createValidatorSecretsLambda(ak: string, sk: string): Promise<any> {
        const url = `/v1/ethereum/validators/aws/lambda/secrets/create`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: ak,
                    secretKey: sk,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create validator secrets lambda function');
            console.error(exc);
            return
        }
    }
    async verifyLambdaFunctionSigner(ak: string, sk: string): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/verify`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: ak,
                    secretKey: sk,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending lambda function verification request');
            console.error(exc);
            return
        }
    }
}
export const awsApiGateway = new AwsApiGateway();

type AuthAWS = {
    region: string;
    accessKey: string;
    secretKey: string;
};

export type AwsRequest = {
    authAWS: AuthAWS;
};
