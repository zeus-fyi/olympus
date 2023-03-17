import {pipe, prop} from 'ramda';
import {getAxiosResponse} from "../../helpers/get-axios-response";
import {authApiGateway} from "../../gateway/login";
import inMemoryJWT from "../../auth/InMemoryJWT";

const sessionIDParse = pipe(getAxiosResponse,prop('sessionID'));
const ttlSeconds = pipe(getAxiosResponse, prop('ttl'));
const userIDParse = pipe(getAxiosResponse, prop('userID'));

const authProvider = {
    login: async (username: string, password: string) =>  {
        try {
            const res = await authApiGateway.sendLoginRequest(username, password);
            console.log(res)
            const statusCode = res.status;
            if (statusCode === 401 || statusCode === 403) {
                inMemoryJWT.ereaseToken();
            }
            if (statusCode === 200 || statusCode === 204) {
                const sessionID = sessionIDParse(res);
                const tokenExpiry = ttlSeconds(res);
                const userID = userIDParse(res);
                inMemoryJWT.setToken(sessionID, tokenExpiry);
                localStorage.setItem("userID", userID);
            }
            console.log(res)
            return res
        } catch (e) {
            console.log(e);
            return e;
        }
    },

    logout: () =>{
        localStorage.removeItem("userID");
        inMemoryJWT.ereaseToken();
    },

    checkAuth: () => {
        return inMemoryJWT.waitForTokenRefresh().then(() => {
            return inMemoryJWT.getToken() ? Promise.resolve() : Promise.reject();
        });
    },

    getPermissions: () => {
        return inMemoryJWT.waitForTokenRefresh().then(() => {
            return inMemoryJWT.getToken() ? Promise.resolve() : Promise.reject();
        });
    },
};

export default authProvider;