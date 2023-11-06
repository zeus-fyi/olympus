import {combineReducers, configureStore} from '@reduxjs/toolkit';
import awsCredentialsReducer from './aws_wizard/aws.wizard.reducer';
import validatorSecretsReducer from './validators/ethereum.validators.reducer';
import clusterBuilderReducer from './clusters/clusters.builder.reducer';
import appsReducer from "./apps/apps.reducer";
import billingReducer from "./billing/billing.reducer";
import resourcesReducer from "./resources/resources.reducer";
import sessionReducer from "./auth/session.reducer";
import loadbalancingReducer from "./loadbalancing/loadbalancing.reducer";
import mevReducer from "./mev/mev.reducer";

const rootReducer = combineReducers({
    sessionState: sessionReducer,
    mev: mevReducer,
    apps: appsReducer,
    loadBalancing: loadbalancingReducer,
    resources: resourcesReducer,
    clusterBuilder: clusterBuilderReducer,
    awsCredentials: awsCredentialsReducer,
    validatorSecrets: validatorSecretsReducer,
    billing: billingReducer,
});

const store = configureStore({
    reducer: rootReducer,
});

export type RootState = ReturnType<typeof rootReducer>;
export type AppDispatch = typeof store.dispatch;

export default store;