import {createSlice} from "@reduxjs/toolkit";
import {MevState} from "./mev.actions";

const initialState: MevState = {
}
const mevSlice = createSlice({
    name: 'mev',
    initialState,
    reducers: {
    }
});

export const {  } = mevSlice.actions;
export default mevSlice.reducer;