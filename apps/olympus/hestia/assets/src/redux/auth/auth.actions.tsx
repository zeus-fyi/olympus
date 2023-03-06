import {pipe, prop} from 'ramda';
import {getAxiosResponse} from "../../helpers/get-axios-response";
import {authApiGateway} from "../../gateway/login";
import inMemoryJWT from "../../auth/InMemoryJWT";
import {AppDispatch} from "../store";

import {LOGOUT,} from "./auth.types";

const tokenParse = pipe(getAxiosResponse,prop('jwtToken'));
const ttlSeconds = pipe(getAxiosResponse, prop('ttl'));
const userIDParse = pipe(getAxiosResponse, prop('userID'));

const authProvider = {
    login: (username: string, password: string) => async () => {
        try {
            const res = await authApiGateway.sendLoginRequest(username, password);
            const statusCode = res.status;
            if (statusCode === 401 || statusCode === 403) {
                inMemoryJWT.ereaseToken();
            }
            if (statusCode === 200) {
                const token = tokenParse(res);
                const tokenExpiry = ttlSeconds(res);
                const userID = userIDParse(res);
                inMemoryJWT.setToken(token, tokenExpiry);
                localStorage.setItem("userID", userID);
            }
            return res
        } catch (e) {
            console.log(e);
            return e;
        }
    },

    logout: async () => async (dispatch: AppDispatch) => {
        const res = await authApiGateway.sendLogoutRequest();
        localStorage.removeItem("user");
        inMemoryJWT.ereaseToken();
        dispatch({
            type: LOGOUT,
        });
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