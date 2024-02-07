import * as React from "react";
import {useEffect} from "react";
import {ClusterSetupContent} from "./ClustersSetup";
import {clustersApiGateway} from "../../../gateway/clusters";
import {useDispatch, useSelector} from "react-redux";
import {setClustersConfigs, updateClusterConfigs} from "../../../redux/clusters/clusters.configs.reducer";
import {RootState} from "../../../redux/store";
import {FormControlLabel, Stack, Switch} from "@mui/material";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";

export default function ClusterConfig() {
    const [loading, setIsLoading] = React.useState(false);

    const dispatch = useDispatch();
    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true);
            try {
                const response = await clustersApiGateway.getExtClustersConfigs();
                dispatch(setClustersConfigs(response.data));
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
    return <ClusterSetupContent loading={loading} setIsLoading={setIsLoading} />;
}

export function ClusterConfigList(props: any) {
    const {loading, setIsLoading} = props;
    const dispatch = useDispatch();
    const clusterConfigs = useSelector((state: RootState) => state.clustersConfigs.clusterConfigs);

    const putExtClusterConfigChanges = async (event: any) => {
        try {
            setIsLoading(true)
            const response = await clustersApiGateway.putExtClustersConfigs(clusterConfigs);
            const statusCode = response.status;
            if (statusCode < 400) {
                // const data = response.data;
            } else {
            }
        } catch (e) {
        } finally {
            setIsLoading(false);
        }
    };

    if (loading) {
        return <div>Loading...</div>;
    }

    const handleChange = (index: number, field: string, value: any) => {
        dispatch(updateClusterConfigs({ index, changes: { [field]: value } }));
    };

    return (
        <div>
            {clusterConfigs && clusterConfigs.map((config, index) => (
                <Stack key={index} direction="row">
                    <Box flexGrow={3} sx={{ mb: 0, mt: 2, mr: 1 }}>
                        <TextField
                            label="Config ID"
                            variant="outlined"
                            value={config.extConfigStrID}
                            InputProps={{
                                readOnly: true,
                            }}
                        />
                    </Box>
                    <Box flexGrow={3} sx={{ mb: 0, mt: 2, mr: 1 }}>
                        <TextField
                            label="Cloud Provider"
                            variant="outlined"
                            value={config.cloudProvider}
                            InputProps={{
                                readOnly: true,
                            }}
                        />
                    </Box>
                    <Box flexGrow={3} sx={{ mb: 0, mt: 2, mr: 1 }}>
                        <TextField
                            label="Region"
                            variant="outlined"
                            value={config.region}
                            InputProps={{
                                readOnly: true,
                            }}
                            onChange={(e) => handleChange(index, 'region', e.target.value)}
                        />
                    </Box>
                    <Box flexGrow={3} sx={{ mb: 0, mt: 2, mr: 1 }}>
                        <TextField
                            label="Context Name"
                            variant="outlined"
                            value={config.context}
                            InputProps={{
                                readOnly: true,
                            }}
                        />
                    </Box>
                    <Box flexGrow={3} sx={{ mb: 0, mt: 2, mr: 1 }}>
                        <TextField
                            label="Context Alias"
                            variant="outlined"
                            value={config.contextAlias}
                            onChange={(e) => handleChange(index, 'contextAlias', e.target.value)}
                        />
                    </Box>
                    <Box flexGrow={3} sx={{ mb: 0, mt: 2, mr: 1 }}>
                        <TextField
                            label="Environment"
                            variant="outlined"
                            value={config.env}
                            onChange={(e) => handleChange(index, 'env', e.target.value)}
                        />
                    </Box>
                    <Box flexGrow={3} sx={{ mb: 0, mt: 2, mr: 1 }}>
                        <FormControlLabel
                            control={
                                <Switch
                                    checked={config.isActive || false}
                                    onChange={(e) => handleChange(index, 'isActive', e.target.checked)}
                                    name="contextToggle"
                                    color="primary"
                                />
                            }
                            label="Context Active"
                        />
                    </Box>
                </Stack>
            ))}
            <Box flexGrow={3} sx={{ mb: 0, mt: 2, mr: 1 }}>
                <Button onClick={(e) => putExtClusterConfigChanges(e)} variant="contained">Update</Button>
            </Box>
        </div>
    );
}
