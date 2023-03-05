import {combineReducers} from 'redux';
import {sessionReducer} from 'redux-react-session';

const reducers = {
    session: sessionReducer,
};

const reducer = combineReducers(reducers);
export default reducer;