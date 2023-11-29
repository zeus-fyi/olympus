import {createSlice, PayloadAction} from "@reduxjs/toolkit";

export interface MenusState {
    openAiPanel: boolean;
    openAppsPanel: boolean;
    openComputePanel: boolean;
}
const initialState: MenusState = {
    openAiPanel: false,
    openAppsPanel: false,
    openComputePanel: false,
}
const menuSlice = createSlice({
    name: 'menus',
    initialState,
    reducers: {
        setOpenAiPanel: (state, action: PayloadAction<boolean>) => {
            state.openAiPanel = action.payload;
        },
        setOpenAppsPanel: (state, action: PayloadAction<boolean>) => {
            state.openAppsPanel = action.payload;
        },
        setOpenComputePanel: (state, action: PayloadAction<boolean>) => {
            state.openComputePanel = action.payload;
        }
    }
});

export const {
    setOpenAiPanel,
    setOpenAppsPanel,
    setOpenComputePanel,
} = menuSlice.actions;
export default menuSlice.reducer