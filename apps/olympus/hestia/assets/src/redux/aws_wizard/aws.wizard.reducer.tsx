import {createSlice, PayloadAction} from "@reduxjs/toolkit";

interface AwsCredentialsState {
    accessKey: string;
    secretKey: string;
    validatorSecretsName: string;
    ageSecretName: string;
    blsSignerLambdaFnUrl: string;
    secretGenLambdaFnUrl: string;
    encKeystoresZipLambdaFnUrl: string;
    depositsGenLambdaFnUrl: string;
    keystoreZip: any
    depositData: [{}]
}
const initialState: AwsCredentialsState = {
    accessKey: '',
    secretKey: '',
    validatorSecretsName: 'mnemonicAndHDWalletEphemery',
    ageSecretName: 'ageEncryptionKeyEphemery',
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
        setValidatorSecretsName: (state, action: PayloadAction<string>) => {
            state.validatorSecretsName = action.payload;
        },
        setAgeSecretName: (state, action: PayloadAction<string>) => {
            state.ageSecretName = action.payload;
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

export const { setAccessKey, setSecretKey, setAgeSecretName, setValidatorSecretsName,
    setBlsSignerLambdaFnUrl, setSecretGenLambdaFnUrl, setEncKeystoresZipLambdaFnUrl, setDepositsGenLambdaFnUrl,
    setKeystoreZip, setDepositData } = awsCredentialsSlice.actions;

export default awsCredentialsSlice.reducer;