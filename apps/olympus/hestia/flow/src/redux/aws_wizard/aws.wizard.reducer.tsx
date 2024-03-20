import {createSlice, PayloadAction} from "@reduxjs/toolkit";

interface AwsCredentialsState {
    accessKey: string;
    secretKey: string;
    validatorSecretsName: string;
    ageSecretName: string;
    externalAccessUserName: string,
    externalAccessSecretName: string,
    blsSignerFunctionName: string;
    blsSignerKeystoresLayerName: string;
    blsSignerLambdaFnUrl: string;
    secretGenLambdaFnUrl: string;
    encKeystoresZipLambdaFnUrl: string;
    depositsGenLambdaFnUrl: string;
    depositData: [{}],
    keystoreLayerNumber: number,
}
export const initialState: AwsCredentialsState = {
    accessKey: '',
    secretKey: '',
    validatorSecretsName: 'mnemonicAndHDWallet',
    ageSecretName: 'ageEncryptionKey',
    externalAccessUserName: 'externalLambdaUser',
    externalAccessSecretName: 'externalZeusLambdaAccessKeys',
    blsSignerFunctionName: 'zeusValidatorsClient',
    blsSignerKeystoresLayerName: 'zeusValidatorsClientKeystores',
    blsSignerLambdaFnUrl: '',
    secretGenLambdaFnUrl: '',
    encKeystoresZipLambdaFnUrl: '',
    depositsGenLambdaFnUrl: '',
    depositData: [{}],
    keystoreLayerNumber: 0,
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
        setSignerFunctionName: (state, action: PayloadAction<string>) => {
            state.blsSignerFunctionName = action.payload;
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
        setDepositData: (state, action: PayloadAction<any>) => {
            state.depositData = action.payload;
        },
        setKeystoreLayerName: (state, action: PayloadAction<string>) => {
            state.blsSignerKeystoresLayerName = action.payload;
        },
        setKeystoreLayerNumber: (state, action: PayloadAction<number>) => {
            state.keystoreLayerNumber = action.payload;
        },
    },
});

export const { setAccessKey, setSecretKey, setAgeSecretName, setValidatorSecretsName,
    setBlsSignerLambdaFnUrl, setSecretGenLambdaFnUrl, setEncKeystoresZipLambdaFnUrl, setDepositsGenLambdaFnUrl,
    setDepositData, setKeystoreLayerNumber, setSignerFunctionName, setKeystoreLayerName } = awsCredentialsSlice.actions;

export default awsCredentialsSlice.reducer;