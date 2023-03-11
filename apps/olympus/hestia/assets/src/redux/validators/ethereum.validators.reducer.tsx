import {createSlice, PayloadAction} from "@reduxjs/toolkit";

interface ValidatorSecretsState {
    hdWalletPw: string;
    mnemonic: string;
}

const initialState: ValidatorSecretsState = {
    hdWalletPw: '',
    mnemonic: '',
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
    },
});

export const { setHdWalletPw, setMnemonic } = validatorSecretsSlice.actions;

export default validatorSecretsSlice.reducer;