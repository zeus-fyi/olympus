import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";
import {AddSkeletonBases} from "./AddSkeletonBases";
import {
    setSelectedComponentBaseName,
    setSelectedSkeletonBaseName
} from "../../../../redux/clusters/clusters.builder.reducer";

export function DefineClusterComponentBaseParams(props: any) {
    const {} = props;
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
                        <SelectedComponentBaseName />
                    </Box>
                    <Box mt={2}>
                        <AddSkeletonBases />
                    </Box>
                </Container>
            </Card>
        </div>
    );
}

export function SelectedComponentBaseName(props: any) {
    const dispatch = useDispatch();
    let cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    let selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const onAccessComponentBase = (selectedComponentBaseName: string) => {
       dispatch(setSelectedComponentBaseName(selectedComponentBaseName));
       const skeletonBaseName = Object.keys(cluster.componentBases[selectedComponentBaseName])[0];
       dispatch(setSelectedSkeletonBaseName(skeletonBaseName));
    };

    let show = Object.keys(cluster.componentBases).length > 0;
    return (
        <div>
            {show &&
            <FormControl variant="outlined" style={{ minWidth: '100%' }}>
                <InputLabel id="network-label">Component Bases</InputLabel>
                <Select
                    labelId="componentBase-label"
                    id="componentBase"
                    value={selectedComponentBaseName}
                    label="Component Base"
                    onChange={(event) => onAccessComponentBase(event.target.value as string)}
                    sx={{ width: '100%' }}
                >
                    {Object.keys(cluster.componentBases).map((key: any, i: number) => (
                        <MenuItem key={i} value={key}>
                            {key}
                        </MenuItem>))
                    }
                </Select>
            </FormControl>
            }
        </div>);
}
