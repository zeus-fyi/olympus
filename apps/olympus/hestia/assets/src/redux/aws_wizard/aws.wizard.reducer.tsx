import {createSlice, PayloadAction} from "@reduxjs/toolkit";

interface AwsCredentialsState {
    accessKey: string;
    secretKey: string;
    agePubKey: string;
    agePrivKey: string;
    blsSignerLambdaFnUrl: string;
    secretGenLambdaFnUrl: string;
    encKeystoresZipLambdaFnUrl: string;
    depositsGenLambdaFnUrl: string;
    keystoreZip: any
    depositData: any
}

const initialState: AwsCredentialsState = {
    accessKey: '',
    secretKey: '',
    agePubKey: '',
    agePrivKey: '',
    blsSignerLambdaFnUrl: '',
    secretGenLambdaFnUrl: '',
    encKeystoresZipLambdaFnUrl: '',
    depositsGenLambdaFnUrl: '',
    keystoreZip: null,
    depositData: [{}]
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
        setBlsSignerLambdaFnUrl: (state, action: PayloadAction<string>) => {
            state.blsSignerLambdaFnUrl = action.payload;
        },
        setSecretGenLambdaFnUrl: (state, action: PayloadAction<string>) => {
            state.secretGenLambdaFnUrl = action.payload;
        },
        setEncKeystoresZipLambdaFnUrl: (state, action: PayloadAction<string>) => {
            state.encKeystoresZipLambdaFnUrl = action.payload;
        },
        setDepositsGenLambdaFnUrl: (state, action: PayloadAction<string>) => {
            state.depositsGenLambdaFnUrl = action.payload;
        },
        setKeystoreZip: (state, action: PayloadAction<any>) => {
            state.keystoreZip = action.payload;
        },
        setDepositData: (state, action: PayloadAction<any>) => {
            state.depositData = action.payload;
        },
    },
});

export const { setAccessKey, setSecretKey, setAgePubKey, setAgePrivKey,
    setBlsSignerLambdaFnUrl, setSecretGenLambdaFnUrl, setEncKeystoresZipLambdaFnUrl, setDepositsGenLambdaFnUrl,
    setKeystoreZip, setDepositData } = awsCredentialsSlice.actions;

export default awsCredentialsSlice.reducer;