import {combineReducers} from 'redux';
import {sessionReducer} from 'redux-react-session';
import clustersReducer from "./clusters/clusters.reducer";

const reducers = {
    clusters: clustersReducer,
    session: sessionReducer,
};

const reducer = combineReducers(reducers);
export default reducer;