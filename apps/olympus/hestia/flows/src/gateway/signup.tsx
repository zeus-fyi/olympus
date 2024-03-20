import {hestiaApi} from './axios/axios';

const config = {
    withCredentials: true,
};

class SignUpApiGateway {
    async sendSignUpRequest(firstName: string, lastName: string, email: string, password: string): Promise<any>  {
        const url = `signup`;
        try {
            return await hestiaApi.post(url, {
                firstName: firstName,
                lastName: lastName,
                email: email,
                password: password,
            }, config)
        } catch (exc) {
            console.error('error sending signup request');
            console.error(exc);
            return
        }
    }
    async verifyEmail(token: string): Promise<any>  {
        const url = `/verify/email/${token}`;
        try {
            return await hestiaApi.get(url, config)
        } catch (exc) {
            console.error('error verifying email request');
            console.error(exc);
            return
        }
    }
    async verifyJWT(token: string): Promise<any>  {
        const url = `/quicknode/dashboard?jwt=${token}`;
        try {
            return await hestiaApi.get(url, config)
        } catch (exc) {
            console.error('error verifying jwt request');
            console.error(exc);
            return
        }
    }
}
export const signUpApiGateway = new SignUpApiGateway();

