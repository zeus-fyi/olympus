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
    async createValidatorsServiceRequest(payload: CreateValidatorServiceRequest): Promise<any>  {
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

// TypeScript interfaces matching the Go types
interface AuthLambdaAWS {
    serviceURL: string;
    secretName: string;
    accessKey: string;
    accessSecret: string;
}

export function createAuthAwsLambda(serviceURL: string, secretName: string, credentials: AwsCredentialIdentity): AuthLambdaAWS {
    return {
        serviceURL: serviceURL,
        secretName: secretName,
        accessKey: credentials.accessKeyId,
        accessSecret: credentials.secretAccessKey,
    };
}

interface ServiceAuthConfig {
    awsAuth: AuthLambdaAWS;
}

interface ServiceRequestWrapper {
    groupName: string;
    protocolNetworkID: number;
    enabled: boolean;
    serviceAuth: ServiceAuthConfig;
}

type ValidatorServiceOrgGroup = {
    pubkey: string;
    feeRecipient: string;
};

type CreateValidatorServiceRequest = {
    serviceRequestWrapper: ServiceRequestWrapper;
    validatorServiceOrgGroupSlice: ValidatorServiceOrgGroup[];
};

export function createValidatorOrgGroup(pubkey: string, feeRecipient: string): ValidatorServiceOrgGroup {
    return {
        pubkey: pubkey,
        feeRecipient: feeRecipient,
    };
}
// Function to create and set the CreateValidatorServiceRequest payload
export function createValidatorServiceRequest(
    keyGroupName: string,
    protocolNetworkID: number,
    externalAwsAuth: AuthLambdaAWS,
    validatorServiceOrgGroups: ValidatorServiceOrgGroup[]
): CreateValidatorServiceRequest {
    const serviceRequestWrapper: ServiceRequestWrapper = {
        groupName: keyGroupName,
        protocolNetworkID: protocolNetworkID,
        enabled: true,
        serviceAuth: {
            awsAuth: {
                serviceURL: externalAwsAuth.serviceURL,
                secretName: externalAwsAuth.secretName,
                accessKey: externalAwsAuth.accessKey,
                accessSecret: externalAwsAuth.accessSecret,
            },
        },
    };

    const hestiaServiceRequest: CreateValidatorServiceRequest = {
        serviceRequestWrapper: serviceRequestWrapper,
        validatorServiceOrgGroupSlice: validatorServiceOrgGroups,
    };

    return hestiaServiceRequest;
}
