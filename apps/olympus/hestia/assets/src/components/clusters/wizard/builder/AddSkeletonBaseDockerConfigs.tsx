import {useDispatch, useSelector} from "react-redux";
import * as React from "react";
import {useEffect} from "react";
import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select, Stack} from "@mui/material";
import {RootState} from "../../../../redux/store";
import Typography from "@mui/material/Typography";
import {DefineDockerParams} from "./DefineDockerImage";
import {
    setSelectedComponentBaseName,
    setSelectedSkeletonBaseName
} from "../../../../redux/clusters/clusters.builder.reducer";
import {AddContainers} from "./AddContainers";

export function AddSkeletonBaseDockerConfigs(props: any) {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const componentBases = cluster.componentBases;
    const componentBaseKeys = Object.keys(componentBases);
    const componentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);

    useEffect(() => {
        dispatch(setSelectedComponentBaseName(componentBaseName));
    }, [dispatch, cluster, componentBaseName]);

    let selectedComponentBaseKey = '';
    if (componentBaseKeys.length > 0) {
        selectedComponentBaseKey = componentBaseKeys[0];
    }

    if (cluster.componentBases === undefined) {
        return <div></div>
    }
    const skeletonBaseKeys = cluster.componentBases[componentBaseName];
    const show = skeletonBaseKeys !== undefined && Object.keys(skeletonBaseKeys).length > 0;
    if (!show) {
        return <div></div>
    }
    return (
        <div>
            <Stack direction="row" spacing={2}>
                <Card sx={{ maxWidth: 800 }}>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Configure Skeleton Base Workloads
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Sets Cluster Skeleton Base Workloads
                        </Typography>
                    </CardContent>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        {show && cluster.componentBases[selectedComponentBaseKey] && Object.keys(skeletonBaseKeys).length > 0 &&
                            <Box mt={2}>
                                <AddContainers />
                            </Box>
                        }
                    </Container>
                </Card>
                {show && cluster.componentBases[selectedComponentBaseKey] && Object.keys(skeletonBaseKeys).length > 0 &&
                    <DefineDockerParams />
                }
            </Stack>
        </div>
    )
}

export function SelectedSkeletonBaseName(props: any) {
    const dispatch = useDispatch();
    const skeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const componentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);

    useEffect(() => {
        dispatch(setSelectedComponentBaseName(componentBaseName));
        dispatch(setSelectedSkeletonBaseName(skeletonBaseName));
    }, [dispatch,skeletonBaseName, componentBaseName]);


    const onAccessSkeletonBase = (selectedSkeletonBaseName: string) => {
        dispatch(setSelectedSkeletonBaseName(selectedSkeletonBaseName));
    };

    if (cluster.componentBases === undefined) {
        return <div></div>
    }

    const skeletonBaseKeys = cluster.componentBases[componentBaseName];
    const show = skeletonBaseKeys !== undefined && Object.keys(skeletonBaseKeys).length > 0;
    if (!show) {
        return <div></div>
    }
    return (
        <FormControl variant="outlined" style={{ minWidth: '100%' }}>
            <InputLabel id="network-label">Skeleton Bases</InputLabel>
            <Select
                labelId="skeletonBase-label"
                id="skeletonBase"
                value={skeletonBaseName}
                label="Skeleton Base"
                onChange={(event) => onAccessSkeletonBase(event.target.value as string)}
                sx={{ width: '100%' }}
            >
                {show && Object.keys(skeletonBaseKeys).map((key: any, i: number) => (
                    <MenuItem key={i} value={key}>
                        {key}
                    </MenuItem>))
                }
            </Select>
        </FormControl>
    );
}