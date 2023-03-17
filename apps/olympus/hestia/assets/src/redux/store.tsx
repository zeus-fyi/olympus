import {combineReducers, configureStore} from '@reduxjs/toolkit';
import awsCredentialsReducer from './aws_wizard/aws.wizard.reducer';
import validatorSecretsReducer from './validators/ethereum.validators.reducer';

const rootReducer = combineReducers({
    awsCredentials: awsCredentialsReducer,
    validatorSecrets: validatorSecretsReducer,
});

const store = configureStore({
    reducer: rootReducer,
});

export type RootState = ReturnType<typeof rootReducer>;
export type AppDispatch = typeof store.dispatch;

export default store;