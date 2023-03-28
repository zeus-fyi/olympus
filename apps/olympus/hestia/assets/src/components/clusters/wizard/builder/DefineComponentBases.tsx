import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";
import {AddSkeletonBases} from "./AddSkeletonBases";
import {
    setSelectedComponentBase,
    setSelectedComponentBaseName
} from "../../../../redux/clusters/clusters.builder.reducer";

export function DefineClusterComponentBaseParams(props: any) {
    const {} = props;
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const componentBaseKeys = Object.keys(cluster.componentBases);

    let selectedKey = '';
    if (componentBaseKeys.length > 0) {
        selectedKey = componentBaseKeys[0];
    }
    let componentBaseObj = cluster.componentBases[selectedKey];
    const onAccessComponentBase = (selectedComponentBaseName: string) => {
        dispatch(setSelectedComponentBaseName(selectedComponentBaseName));
        componentBaseObj = cluster.componentBases[selectedKey]
        dispatch(setSelectedComponentBase(componentBaseObj));

    };

    // TODO, when component base is changed needs to update the skeleton base list

    return (
        <div>
            <Card sx={{ maxWidth: 500 }}>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Define Skeleton Base Elements for Component Bases
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Sets Skeleton Base Elements for Component Bases
                    </Typography>
                </CardContent>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Box mt={2}>
                        <SelectedComponentBaseName componentBaseKeys={componentBaseKeys} onAccessComponentBase={onAccessComponentBase} />
                    </Box>
                    <Box mt={2}>
                        <AddSkeletonBases componentBase={componentBaseObj} />
                    </Box>
                </Container>
            </Card>
        </div>
    );
}

export function SelectedComponentBaseName(props: any) {
    const {onAccessComponentBase} = props;
    const componentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const componentBaseKeys = Object.keys(cluster.componentBases);
    return (
        <FormControl variant="outlined" style={{ minWidth: '100%' }}>
            <InputLabel id="network-label">Component Bases</InputLabel>
            <Select
                labelId="componentBase-label"
                id="componentBase"
                value={componentBaseName}
                label="Component Base"
                onChange={(event) => onAccessComponentBase(event.target.value as string)}
                sx={{ width: '100%' }}
            >
                {componentBaseKeys.map((key: any) => (
                    <MenuItem key={key} value={key}>
                        {key}
                    </MenuItem>))
                }
            </Select>
        </FormControl>
    );
}