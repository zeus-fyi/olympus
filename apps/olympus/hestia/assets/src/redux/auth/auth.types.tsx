export const REGISTER_SUCCESS = "REGISTER_SUCCESS";
export const REGISTER_FAIL = "REGISTER_FAIL";
export const LOGIN_SUCCESS = "LOGIN_SUCCESS";
export const LOGIN_FAIL = "LOGIN_FAIL";
export const LOGOUT = "LOGOUT";

export const SET_MESSAGE = "SET_MESSAGE";
export const CLEAR_MESSAGE = "CLEAR_MESSAGE";

export interface User {
    userID: string
    email: string
}

export interface LoginResponse {
    userID: number
    sessionID: string
    ttl: number
    isBillingSetup: boolean
    isInternal: boolean
}

export interface LoginRequest {
    email: string
    password: string
}