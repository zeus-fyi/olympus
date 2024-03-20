import {AwsCredentialIdentity} from "@aws-sdk/types/dist-types/identity";
import {AwsClient} from 'aws4fetch'

class AwsLambdaApiGateway {
    async invokeValidatorSecretsGeneration(url:string, credentials: AwsCredentialIdentity, mnemonicHdPwSecretName: string, ageKeySecretName: string): Promise<any> {
        try {
            const payload = {
                mnemonicAndHDWalletSecretName: mnemonicHdPwSecretName,
                ageSecretName: ageKeySecretName,
            }
            const headers = {
                "Content-Type": "application/json",
            };
            const accessKeyId = credentials.accessKeyId;
            const secretAccessKey = credentials.secretAccessKey;
            const service = "lambda";
            const region = "us-west-1";
            const body = JSON.stringify(payload);
            const aws = new AwsClient({
                accessKeyId,
                secretAccessKey,
                service,
                region,
            })
            const request = new Request(
                `${url}`,
                {
                    method: "POST",
                    headers,
                    body,
                }
            );
            return await aws.fetch(request)
        } catch (exc) {
            console.error('error sending create new secrets gen via lambda');
            console.error(exc);
            return
        }
    }
    async invokeValidatorDepositsGeneration(url: string, credentials: AwsCredentialIdentity, network: string,
                                            mnemonicHdPwSecretName: string, validatorCount: number,
                                            hdOffset: number, wc: string): Promise<any> {
        try {
            let payload: DepositGenerationPayload = {
                mnemonicAndHDWalletSecretName: mnemonicHdPwSecretName,
                validatorCount: validatorCount,
                hdOffset: hdOffset,
                network: network,
            }
            if (network === 'Goerli') {
                const forkVersionBytes = new Uint8Array([0x00, 0x00, 0x10, 0x20]);
                payload.forkVersion =  Array.from(forkVersionBytes)
            }
            if (network === 'Mainnet') {
                const forkVersionBytes = new Uint8Array([0x00, 0x00, 0x00, 0x00]);
                payload.forkVersion =  Array.from(forkVersionBytes)
            }
            if (wc.length > 0) {
                payload.withdrawalAddress = wc;
            }
            const headers = {
                "Content-Type": "application/json",
            };
            const accessKeyId = credentials.accessKeyId;
            const secretAccessKey = credentials.secretAccessKey;
            const service = "lambda";
            const region = "us-west-1";
            const body = JSON.stringify(payload);
            const aws = new AwsClient({
                accessKeyId,
                secretAccessKey,
                service,
                region,
            })
            const request = new Request(
                `${url}`,
                {
                    method: "POST",
                    headers,
                    body,
                }
            );
            return await aws.fetch(request)
        } catch (exc) {
            console.error('error sending create validator deposits to lambda function');
            console.error(exc);
            return
        }
    }
    async invokeEncryptedKeystoresZipGeneration(url: string, credentials: AwsCredentialIdentity, ageSecretName: string, mnemonicHdPwSecretName: string, validatorCount: number, hdOffset: number): Promise<any> {
        try {
            const payload = {
                ageSecretName: ageSecretName,
                mnemonicAndHDWalletSecretName: mnemonicHdPwSecretName,
                validatorCount: validatorCount,
                hdOffset: hdOffset,
            }
            const headers = {
                "Content-Type": "application/json",
            };
            const accessKeyId = credentials.accessKeyId;
            const secretAccessKey = credentials.secretAccessKey;
            const service = "lambda";
            const region = "us-west-1";
            const body = JSON.stringify(payload);
            const aws = new AwsClient({
                accessKeyId,
                secretAccessKey,
                service,
                region,
            })
            const request = new Request(
                `${url}`,
                {
                    method: "POST",
                    headers,
                    body,
                }
            );
            return await aws.fetch(request)
        } catch (exc) {
            console.error('error sending create validator encrypted keystores zip request to lambda function');
            console.error(exc);
            return
        }
    }
}

interface DepositGenerationPayload {
    mnemonicAndHDWalletSecretName: string;
    validatorCount: number;
    hdOffset: number;
    network: string;
    withdrawalAddress?: string;
    forkVersion?: any; // Change the type of forkVersion to
}

export const awsLambdaApiGateway = new AwsLambdaApiGateway();
