import {createSlice, PayloadAction} from "@reduxjs/toolkit";

export interface SessionState {
    sessionAuth: boolean;
    isInternal: boolean;
}

const initialState: SessionState = {
    sessionAuth: false,
    isInternal: false,
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
    }
});

export const { setSessionAuth, setInternalAuth } = sessionStateSlice.actions;
export default sessionStateSlice.reducer