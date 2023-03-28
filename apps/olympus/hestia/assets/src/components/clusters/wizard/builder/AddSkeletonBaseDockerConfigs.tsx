import {useDispatch, useSelector} from "react-redux";
import * as React from "react";
import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import {RootState} from "../../../../redux/store";
import Typography from "@mui/material/Typography";
import {DefineDockerParams} from "./DefineDockerImage";
import {setSelectedSkeletonBaseName} from "../../../../redux/clusters/clusters.builder.reducer";

export function AddSkeletonBaseDockerConfigs(props: any) {
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const componentBases = cluster.componentBases;
    const componentBaseKeys = Object.keys(componentBases);
    const componentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);

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
            <Card sx={{ maxWidth: 1000 }}>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Configure Skeleton Base Workloads
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Sets Cluster Skeleton Base Workloads
                    </Typography>
                </CardContent>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    { show && cluster.componentBases[selectedComponentBaseKey] && Object.keys(skeletonBaseKeys).length > 0 &&
                        <Box mt={2}>
                            <SelectedSkeletonBaseName />
                        </Box>
                    }
                </Container>
                {show && cluster.componentBases[selectedComponentBaseKey] && Object.keys(skeletonBaseKeys).length > 0 &&
                    <DefineDockerParams />
                }
            </Card>
        </div>
    )
}

export function SelectedSkeletonBaseName(props: any) {
    const dispatch = useDispatch();
    const skeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const onAccessSkeletonBase = (selectedSkeletonBaseName: string) => {
        dispatch(setSelectedSkeletonBaseName(selectedSkeletonBaseName));
    };
    const componentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);

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