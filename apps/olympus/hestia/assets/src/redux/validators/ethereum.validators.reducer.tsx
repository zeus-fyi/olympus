import {createSlice, PayloadAction} from "@reduxjs/toolkit";

interface ValidatorSecretsState {
    hdWalletPw: string;
    mnemonic: string;
    hdOffset: number;
    validatorCount: number;
}

const initialState: ValidatorSecretsState = {
    hdWalletPw: '',
    mnemonic: '',
    hdOffset: 0,
    validatorCount: 1
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
    },
});

export const { setHdWalletPw, setMnemonic } = validatorSecretsSlice.actions;

export default validatorSecretsSlice.reducer;