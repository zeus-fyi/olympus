import inMemoryJWT from "../auth/InMemoryJWT";
import {hestiaApi} from "./axios/axios";

class MevApiGateway {
    async getDashboardInfo(): Promise<any>  {
        const url = `/web/internal/v1/mev/dashboard`;
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
            console.error('error sending mev dashboard request');
            console.error(exc);
            return
        }
    }
}
export const mevApiGateway = new MevApiGateway();
