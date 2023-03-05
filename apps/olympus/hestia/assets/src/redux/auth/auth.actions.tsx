import {pipe, prop} from 'ramda';
import {getAxiosResponse} from "../../helpers/get-axios-response";
import {authApiGateway} from "../../gateway/login";
import inMemoryJWT from "../../auth/InMemoryJWT";
import {AppDispatch} from "../store";

export const REGISTER_SUCCESS = "REGISTER_SUCCESS";
export const REGISTER_FAIL = "REGISTER_FAIL";
export const LOGIN_SUCCESS = "LOGIN_SUCCESS";
export const LOGIN_FAIL = "LOGIN_FAIL";
export const LOGOUT = "LOGOUT";

const jwtTokenParse = pipe(getAxiosResponse,prop('jwtToken'));
const ttlSeconds = pipe(getAxiosResponse, prop('ttl'));

const authProvider = {
    login: async (username: string, password: string) => {
        try {
            const res = await authApiGateway.sendLoginRequest(username, password);
            const statusCode = res.status;
            if (statusCode === 401 || statusCode === 403) {
                inMemoryJWT.ereaseToken();
            } else {
                const token = jwtTokenParse(res);
                const tokenExpiry = ttlSeconds(res);
                inMemoryJWT.setToken(token, tokenExpiry);
            }
            return res
        } catch (e) {
            console.log(e);
            return e;
        }
    },

    logout: async () => async (dispatch: AppDispatch) => {
        const res = await authApiGateway.sendLogoutRequest();
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