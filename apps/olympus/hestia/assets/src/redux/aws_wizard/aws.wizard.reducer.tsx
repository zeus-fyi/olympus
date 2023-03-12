import {createSlice, PayloadAction} from "@reduxjs/toolkit";

interface AwsCredentialsState {
    accessKey: string;
    secretKey: string;
    agePubKey: string;
    agePrivKey: string;
}

const initialState: AwsCredentialsState = {
    accessKey: '',
    secretKey: '',
    agePubKey: '',
    agePrivKey: '',
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
        setAgePubKey: (state, action: PayloadAction<string>) => {
            state.agePubKey = action.payload;
        },
        setAgePrivKey: (state, action: PayloadAction<string>) => {
            state.agePrivKey = action.payload;
        },
    },
});

export const { setAccessKey, setSecretKey, setAgePubKey, setAgePrivKey } = awsCredentialsSlice.actions;

export default awsCredentialsSlice.reducer;