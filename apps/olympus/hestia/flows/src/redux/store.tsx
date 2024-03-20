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
import aiReducer from "./ai/ai.reducer";
import menusReducer from "./menus/menus.reducer";
import clustersConfigsReducer from "./clusters/clusters.configs.reducer";
import flowsReducer from "./flows/flows.reducer";

const rootReducer = combineReducers({
    flows: flowsReducer,
    sessionState: sessionReducer,
    mev: mevReducer,
    apps: appsReducer,
    loadBalancing: loadbalancingReducer,
    resources: resourcesReducer,
    clusterBuilder: clusterBuilderReducer,
    awsCredentials: awsCredentialsReducer,
    validatorSecrets: validatorSecretsReducer,
    billing: billingReducer,
    ai: aiReducer,
    menus: menusReducer,
    clustersConfigs: clustersConfigsReducer,
});

const store = configureStore({
    reducer: rootReducer,
});

export type RootState = ReturnType<typeof rootReducer>;
export type AppDispatch = typeof store.dispatch;

export default store;