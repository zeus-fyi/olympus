import axios from './axios/axios';

const config = {
};

class AuthApiGateway {
    async sendLoginRequest(email: string, password: string): Promise<any>  {
        const url = `login`;
        try {
            return await axios.post(url, {
                email: email,
                password: password,
            }, config)
        } catch (exc) {
            console.error('error sending login request');
            console.error(exc);
            return
        }
    }

    async sendLogoutRequest()  {
        const url = `logout`;
        try {
            return await axios.get(url, config)
        } catch (exc) {
            console.error('error sending login request');
            console.error(exc);
            return
        }
    }
}
export const authApiGateway = new AuthApiGateway();

