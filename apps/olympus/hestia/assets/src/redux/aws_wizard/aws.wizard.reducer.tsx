import {createSlice, PayloadAction} from "@reduxjs/toolkit";

interface AwsCredentialsState {
    accessKey: string;
    secretKey: string;
}

const initialState: AwsCredentialsState = {
    accessKey: '',
    secretKey: '',
};

const awsCredentialsSlice = createSlice({
    name: 'awsCredentials',
    initialState,
    reducers: {
        setAccessKey: (state, action: PayloadAction<string>) => {
            state.accessKey = action.payload;
        },
        setSecretKey: (state, action: PayloadAction<string>) => {
            state.secretKey = action.payload;
        },
    },
});

export const { setAccessKey, setSecretKey } = awsCredentialsSlice.actions;

export default awsCredentialsSlice.reducer;