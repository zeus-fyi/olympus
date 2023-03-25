import {createSlice, PayloadAction} from "@reduxjs/toolkit";

interface ValidatorSecretsState {
    hdWalletPw: string;
    mnemonic: string;
    hdOffset: number;
    validatorCount: number;
    network: string;
    feeRecipient: string;
    keyGroupName: string;
    networkAppended : boolean;
    authorizedNetworks: [string]
}

const initialState: ValidatorSecretsState = {
    hdWalletPw: '',
    mnemonic: '',
    hdOffset: 0,
    validatorCount: 1,
    network: 'Ephemery',
    feeRecipient: '',
    keyGroupName: 'DemoKeyGroup',
    networkAppended : false,
    authorizedNetworks: ['Ephemery'],
};

const validatorSecretsSlice = createSlice({
    name: 'validatorSecrets',
    initialState,
    reducers: {
        setHdWalletPw: (state, action: PayloadAction<string>) => {
            state.hdWalletPw = action.payload;
        },
        setMnemonic: (state, action: PayloadAction<string>) => {
            state.mnemonic = action.payload;
        },
        setValidatorCount: (state, action: PayloadAction<number>) => {
            state.validatorCount = action.payload;
        },
        setHdOffset: (state, action: PayloadAction<number>) => {
            state.hdOffset = action.payload;
        },
        setNetworkName: (state, action: PayloadAction<string>) => {
            state.network = action.payload;
        },
        setFeeRecipient: (state, action: PayloadAction<string>) => {
            state.feeRecipient = action.payload;
        },
        setKeyGroupName: (state, action: PayloadAction<string>) => {
            state.keyGroupName = action.payload;
        },
        setNetworkAppended: (state, action: PayloadAction<boolean>) => {
            state.networkAppended = action.payload;
        },
        setAuthorizedNetworks: (state, action: PayloadAction<[string]>) => {
            state.authorizedNetworks = action.payload;
        },
    },
});

export const { setHdOffset, setValidatorCount, setNetworkName, setFeeRecipient, setKeyGroupName, setNetworkAppended, setAuthorizedNetworks} = validatorSecretsSlice.actions;

export default validatorSecretsSlice.reducer;