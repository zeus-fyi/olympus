import {combineReducers, configureStore} from '@reduxjs/toolkit';
import {sessionReducer, sessionService, SessionServiceOptions} from 'redux-react-session';
import awsCredentialsReducer from './aws_wizard/aws.wizard.reducer';

const rootReducer = combineReducers({
    session: sessionReducer,
    awsCredentials: awsCredentialsReducer,
});

const store = configureStore({
    reducer: rootReducer,
});

const validateSession = (session: any) => {
    // check if your session is still valid
    return true;
};

const options: SessionServiceOptions = {
    refreshOnCheckAuth: true,
    driver: 'COOKIES',
    validateSession,
    expires: 3600,
};

sessionService.initSessionService(store, options)
    .then(() => console.log('Redux React Session is ready and a session was refreshed from your storage'))
    .catch(() => console.log('Redux React Session is ready and there is no session in your storage'));

export type RootState = ReturnType<typeof rootReducer>;
export type AppDispatch = typeof store.dispatch;

export default store;