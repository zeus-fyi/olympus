import * as React from "react";
import {useEffect} from "react";
import {useDispatch} from "react-redux";
import {clustersApiGateway} from "../../../gateway/clusters";
import {setAuthedClustersConfigs} from "../../../redux/clusters/clusters.configs.reducer";

export default function CloudProviderResources() {
    const [loading, setIsLoading] = React.useState(false);

    const dispatch = useDispatch();
    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true);
            try {
                const response = await clustersApiGateway.getPrivateAuthedClustersConfigs();
                dispatch(setAuthedClustersConfigs(response.data));
            } catch (error) {
                console.log("error", error);
            } finally {
                setIsLoading(false);
            }}
        fetchData().then(r => '');
    }, []);

    if (loading) {
        return <div>Loading...</div>;
    }
    return <div>CloudProviderResources</div>;
}