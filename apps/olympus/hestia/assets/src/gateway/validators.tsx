import {hestiaApi} from './axios/axios';

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
    async generateValidatorsDepositData(mnemonicPhrase: string, hdWalletPw: string, count: number, offset: number): Promise<any>  {
        const url = `/v1/ethereum/validators/aws/generation`;
        try {
            const sessionID = localStorage.getItem("sessionID");
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`,
                }}
            const payload = {
                agePubKey: '',
                agePrivKey: '',
                mnemonic: mnemonicPhrase,
                hdWalletPw: hdWalletPw,
                validatorCount: count,
                hdOffset: offset,
            }
            return await hestiaApi.post(url, payload, config)
        } catch (exc) {
            console.error('error sending create validator deposits');
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

