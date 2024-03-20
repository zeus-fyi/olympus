import {clustersApiGateway} from "../../gateway/clusters";

export const SET_CLUSTERS = 'SET_CLUSTERS';
export const GET_CLUSTERS_FAIL = 'GET_CLUSTERS_FAIL';

export const setClusters = (clusters: any) => ({ type: SET_CLUSTERS, payload: clusters });
export const getClustersFail = (clusters: any) => ({ type: GET_CLUSTERS_FAIL, payload: null});

export const fetchClusters = (userID: number) => async (dispatch: any) => {
    if (!userID) {
        return;
    }
    try {
        const res = await clustersApiGateway.getClusters();
        dispatch(setClusters(res.data));
    } catch (exc) {
        dispatch(getClustersFail(null));
        console.error('error while loading users clusters');
        console.error(exc);
    }
};
