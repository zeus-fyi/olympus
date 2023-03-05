import {LOGIN_SUCCESS, REGISTER_SUCCESS,} from "./auth.actions";


export const reducer = (state = {}, action: { type: any; }) => {
    switch (action.type) {
        case LOGIN_SUCCESS:
            return { ...state, loggedIn: true };
        case REGISTER_SUCCESS:
            return { ...state, registered: true };
        default:
            return state;
    }
};
