import {createSlice, PayloadAction} from "@reduxjs/toolkit";

export interface SessionState {
    sessionAuth: boolean;
    isInternal: boolean;
    isBillingSetup: boolean;
}

const initialState: SessionState = {
    sessionAuth: false,
    isInternal: false,
    isBillingSetup: false,
}

const sessionStateSlice = createSlice({
    name: 'sessionState',
    initialState,
    reducers: {
        setSessionAuth: (state, action: PayloadAction<boolean>) => {
            state.sessionAuth = action.payload;
        },
        setInternalAuth: (state, action: PayloadAction<boolean>) => {
            state.isInternal = action.payload;
        },
        setIsBillingSetup: (state, action: PayloadAction<boolean>) => {
            state.isBillingSetup = action.payload;
        },
    }
});

export const {
    setSessionAuth,
    setInternalAuth,
    setIsBillingSetup,
} = sessionStateSlice.actions;
export default sessionStateSlice.reducer