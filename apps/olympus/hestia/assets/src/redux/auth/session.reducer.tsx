import {createSlice, PayloadAction} from "@reduxjs/toolkit";

export interface SessionState {
    sessionAuth: boolean;
}

const initialState: SessionState = {
    sessionAuth: false,
}
const sessionStateSlice = createSlice({
    name: 'sessionState',
    initialState,
    reducers: {
        setSessionAuth: (state, action: PayloadAction<boolean>) => {
            state.sessionAuth = action.payload;
        },
    }
});

export const { setSessionAuth } = sessionStateSlice.actions;
export default sessionStateSlice.reducer