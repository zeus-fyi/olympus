import {artemisApi, hestiaApi} from './axios/axios';
import {AwsCredentialIdentity} from "@aws-sdk/types/dist-types/identity";
import inMemoryJWT from "../auth/InMemoryJWT";

class ValidatorsApiGateway {
    async getValidators(): Promise<any>  {
        const url = `/v1/ethereum/validators/service/info`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
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
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create lambda function keystores layer');
            console.error(exc);
            return
        }
    }
    async depositValidatorsServiceRequest(payload: CreateValidatorsDepositServiceRequest): Promise<any>  {
        const url = `/v1/ethereum/validators/create`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await artemisApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create lambda function keystores layer');
            console.error(exc);
            return
        }
    }
    async getAuthedValidatorsServiceRequest(): Promise<any>  {
        const url = `/v1/users/services`;
        try {
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.get(url, config)
        } catch (exc) {
            console.error('error sending getAuthedValidatorsServiceRequest');
            console.error(exc);
            return
        }
    }
}
export const validatorsApiGateway = new ValidatorsApiGateway();

export interface ValidatorDepositDataJSON {
    pubkey: string;
    withdrawal_credentials: string;
    signature: string;
    deposit_data_root: string;
    amount: number;
    deposit_message_root: string;
    fork_version: string;
}

export interface ValidatorDepositDataRxJSON {
    pubkey: string;
    rx: string;
}

export function createValidatorsDepositsDataJSON(
    pubkey: string,
    withdrawalCredentials: string,
    signature: string,
    depositDataRoot: string,
    amount: number,
    depositMessageRoot: string,
    forkVersion: string,
): ValidatorDepositDataJSON {
    return {
        pubkey,
        withdrawal_credentials: withdrawalCredentials,
        signature,
        amount,
        deposit_data_root: depositDataRoot,
        deposit_message_root: depositMessageRoot,
        fork_version: forkVersion,
    };
}

type CreateValidatorsDepositServiceRequest = {
    network: string;
    validatorDepositSlice: ValidatorDepositDataJSON[];
};

export function createValidatorsDepositServiceRequest(network: string, validatorServiceOrgGroupSlice: ValidatorDepositDataJSON[]): CreateValidatorsDepositServiceRequest {
    return {
        network: network,
        validatorDepositSlice: validatorServiceOrgGroupSlice,
    }
}

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

interface ValidatorDepositParams {
    pubkey: string;
    withdrawal_credentials: string;
    signature: string;
    deposit_data_root: string;
}

interface ExtendedDepositParams extends ValidatorDepositParams {
    amount: number;
    deposit_message_root: string;
}