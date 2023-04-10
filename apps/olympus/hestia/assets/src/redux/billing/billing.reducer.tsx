import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {BillingState} from "./billing.types";

const initialState: BillingState = {
    stripeCustomerID: '',
}
const billingSlice = createSlice({
    name: 'billing',
    initialState,
    reducers: {
        setStripeCustomerID: (state, action: PayloadAction<string>) => {
            state.stripeCustomerID = action.payload;
        },
    }
});

export const { setStripeCustomerID } = billingSlice.actions;
export default billingSlice.reducer