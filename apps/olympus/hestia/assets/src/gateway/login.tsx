import axios from './axios/axios';

class LoginApiGateway {
    async sendLoginRequest(email: string, password: string) {
        const url = `login`;
        await axios.post(url, {
            email: email,
            password: password,
        });
    }
}
export const loginApiGateway = new LoginApiGateway();