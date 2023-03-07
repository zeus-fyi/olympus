import reducer from './reducer';
import {configureStore} from '@reduxjs/toolkit'
import {sessionService, SessionServiceOptions} from 'redux-react-session';
import {Store} from 'redux';
import {setupListeners} from "@reduxjs/toolkit/query";

export const store = configureStore({
    reducer: {
        reducer,
    },
})
const validateSession = (session: any) => {
    // check if your session is still valid
    return true;
}
const options: SessionServiceOptions = { refreshOnCheckAuth: true, driver: 'COOKIES', validateSession, expires: 3600 };

sessionService.initSessionService(store as Store, options)
    .then(() => console.log('Redux React Session is ready and a session was refreshed from your storage'))
    .catch(() => console.log('Redux React Session is ready and there is no session in your storage'));

export default store;
// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch

setupListeners(store.dispatch)
