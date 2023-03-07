
export const SET_CLUSTERS = 'SET_CLUSTERS';

export const setClusters = (clusters: any) => ({ type: SET_CLUSTERS, payload: clusters });

export const fetchClusters = (user: any) => async (dispatch: any) => {
    if (!user) {
        return;
    }
    try {

    } catch (exc) {
        console.error('error while loading users clusters');
        console.error(exc);
    }
};
