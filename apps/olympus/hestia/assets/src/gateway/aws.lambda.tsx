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
            console.error('error sending create new internal lambda user');
            console.error(exc);
            return
        }
    }
    async invokeValidatorDepositsGeneration(url:string, credentials: AwsCredentialIdentity,network: string, mnemonicHdPwSecretName: string, validatorCount: number, hdOffset: number): Promise<any> {
        try {
            const payload = {
                validatorCount: validatorCount,
                hdOffset: hdOffset,
                network: network,
            }
            let headers = {
                "Content-Type": "application/json",
            };
            let accessKeyId = credentials.accessKeyId;
            let secretAccessKey = credentials.secretAccessKey;
            let service = "lambda";
            let region = "us-west-1";
            let body = JSON.stringify(payload);
            let aws = new AwsClient({
                accessKeyId,
                secretAccessKey,
                service,
                region,
            })
            let request = new Request(
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
}

export const awsLambdaApiGateway = new AwsLambdaApiGateway();
