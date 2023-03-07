import {GET_CLUSTERS_FAIL, SET_CLUSTERS} from './clusters.actions';

const getClusterItems: any = localStorage.getItem("clusters")
const clusters: any = JSON.parse(getClusterItems);

const initialState = clusters
    ? { hasClusters: true, clusters }
    : { hasClusters: false, clusters: null };

export default function clustersReducer(state = initialState, action: any ) {
    const { type, payload } = action;

    switch (type) {
        case SET_CLUSTERS:
            return {
                ...state,
                hasClusters: true,
                clusters: payload.clusters,
            };
        case GET_CLUSTERS_FAIL:
            return {
                ...state,
                hasClusters: true,
                clusters: payload.clusters,
            };
        default:
            return state;
    }
}

export const getClusters = (state: any) => state.clusters;