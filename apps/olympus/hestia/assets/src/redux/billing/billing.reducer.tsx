import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {BillingState} from "./billing.types";

const initialState: BillingState = {
    stripeCustomerClientSecret: '',
}
const billingSlice = createSlice({
    name: 'billing',
    initialState,
    reducers: {
        setStripeCustomerClientSecret: (state, action: PayloadAction<string>) => {
            state.stripeCustomerClientSecret = action.payload;
        },
    }
});

export const { setStripeCustomerClientSecret } = billingSlice.actions;
export default billingSlice.reducer