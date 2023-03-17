import {hestiaApi} from './axios/axios';
import {AwsCredentialIdentity} from "@aws-sdk/types/dist-types/identity";
import inMemoryJWT from "../auth/InMemoryJWT";

class AwsApiGateway {
    async createInternalLambdaUser(credentials: AwsCredentialIdentity): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/user/internal/lambda/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: credentials.accessKeyId,
                    secretKey: credentials.secretAccessKey,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create new internal lambda user');
            console.error(exc);
            return
        }
    }
    async createExternalLambdaUser(credentials: AwsCredentialIdentity): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/user/external/lambda/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: credentials.accessKeyId,
                    secretKey: credentials.secretAccessKey,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create new external lambda user');
            console.error(exc);
            return
        }
    }
    async getLambdaFunctionURL(credentials: AwsCredentialIdentity, functionName: string): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/url`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload: AwsRequestSignerCreationRequest = {
                functionName: functionName,
                keystoresLayerName: "",
                authAWS: {
                    region: "us-west-1",
                    accessKey: credentials.accessKeyId,
                    secretKey: credentials.secretAccessKey,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create new external lambda user');
            console.error(exc);
            return
        }
    }
    async createOrFetchExternalLambdaUserAccessKeys(credentials: AwsCredentialIdentity, externalUserName: string, externalAccessSecretName: string): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/external/user/access/keys/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload: AwsRequestSignerExternalUserAccessCreationRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: credentials.accessKeyId,
                    secretKey: credentials.secretAccessKey,
                },
                externalUserName: externalUserName,
                externalAccessSecretName: externalAccessSecretName,
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create new external lambda user');
            console.error(exc);
            return
        }
    }
    async createLambdaSignerFunction(credentials: AwsCredentialIdentity, functionName: string, keystoresLayerName: string): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/signer/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload: AwsRequestSignerCreationRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: credentials.accessKeyId,
                    secretKey: credentials.secretAccessKey,
                },
                functionName: functionName,
                keystoresLayerName: keystoresLayerName,
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create lambda function');
            console.error(exc);
            return
        }
    }
    async createLambdaFunctionKeystoresLayer(credentials: AwsCredentialIdentity, keystoresLayerName: string, keystoresZip: Blob): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/signer/keystores/layer/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const formData = new FormData(); // Create a new FormData object
            const zipFile = new File([keystoresZip], 'keystores.zip', { type: 'application/zip' }); // Create a new zip file
            formData.append('authAWS', JSON.stringify({
                region: "us-west-1",
                accessKey: credentials.accessKeyId,
                secretKey: credentials.secretAccessKey,
            })); // Add authAWS as a stringified JSON value
            formData.append('keystoresLayerName', keystoresLayerName); // Add keystoresLayerName
            formData.append('keystoresZip', zipFile); // Add keystoresZip as a blob
            return await hestiaApi.post(url, formData, config)
        } catch (exc) {
            console.error('error sending create lambda keystores layer');
            console.error(exc);
            return
        }
    }
    async createValidatorsDepositDataLambda(credentials: AwsCredentialIdentity): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/deposits/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: credentials.accessKeyId,
                    secretKey: credentials.secretAccessKey,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create validator deposits lambda function');
            console.error(exc);
            return
        }
    }
    async createValidatorsAgeEncryptedKeystoresZipLambda(credentials: AwsCredentialIdentity): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/keystores/zip/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: credentials.accessKeyId,
                    secretKey: credentials.secretAccessKey,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create validator encrypted keystores lambda function');
            console.error(exc);
            return
        }
    }
    async createValidatorSecretsLambda(credentials: AwsCredentialIdentity): Promise<any> {
        const url = `/v1/ethereum/validators/aws/lambda/secrets/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const payload: AwsRequest = {
                authAWS: {
                    region: "us-west-1",
                    accessKey: credentials.accessKeyId,
                    secretKey: credentials.secretAccessKey,
                },
            };
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create validator secrets lambda function');
            console.error(exc);
            return
        }
    }
    async verifyLambdaFunctionSigner(credentials: AwsCredentialIdentity, ageSecretName: string, fnUrl: string, depositSlice: [{}]): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/lambda/verify`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            const keySlice = depositSlice.map((v: any) =>
                v.pubkey
            );
            const payload: AwsVerifyLambdaSignerRequest = {
                keySlice: keySlice,
                functionURL: fnUrl,
                secretName: ageSecretName,
                authAWS: {
                    region: "us-west-1",
                    accessKey: credentials.accessKeyId,
                    secretKey: credentials.secretAccessKey,
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

export type AwsVerifyLambdaSignerRequest = {
    authAWS: AuthAWS;
    keySlice: any;
    secretName: string;
    functionURL: string;
};

export type AwsRequestSignerExternalUserAccessCreationRequest = {
    authAWS: AuthAWS;
    externalUserName: string;
    externalAccessSecretName: string;
};

export type AwsRequestSignerCreationRequest = {
    authAWS: AuthAWS;
    functionName: string;
    keystoresLayerName: string;
};

export type AwsRequestKeystoreLayerCreationRequest = {
    authAWS: AuthAWS;
    keystoresLayerName: string;
    keystoresZip: Blob;
};

type AuthAWS = {
    region: string;
    accessKey: string;
    secretKey: string;
};

export type AwsRequest = {
    authAWS: AuthAWS;
};

function createKeySlice(
    pubkey: string,
) {
    return {pubkey};
}