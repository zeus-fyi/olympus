import {hestiaApi} from './axios/axios';
import inMemoryJWT from "../auth/InMemoryJWT";

class StripeApiGateway {
    async getClientSecret(): Promise<any>  {
        const url = `/v1/stripe/customer/id`;
            const sessionID = inMemoryJWT.getToken();
            let config = {
                headers: {
                    'Authorization': `Bearer ${sessionID}`
                },
                withCredentials: true,
            }
            return await hestiaApi.get(url, config)
    }
}
export const stripeApiGateway = new StripeApiGateway();

