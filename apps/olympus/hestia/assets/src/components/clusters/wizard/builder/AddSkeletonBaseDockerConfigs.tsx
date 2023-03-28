import {useDispatch, useSelector} from "react-redux";
import * as React from "react";
import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import {RootState} from "../../../../redux/store";
import {SelectedComponentBaseName} from "./DefineComponentBases";
import Typography from "@mui/material/Typography";
import {DefineDockerParams} from "./DefineDockerImage";

export function AddSkeletonBaseDockerConfigs(props: any) {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const componentBases = cluster.componentBases;
    const componentBaseKeys = Object.keys(componentBases);

    let selectedComponentBaseKey = '';
    if (componentBaseKeys.length > 0) {
        selectedComponentBaseKey = componentBaseKeys[0];
    }
    const [componentBase, setComponentBase] = React.useState(selectedComponentBaseKey);
    const onAccessComponentBase = (selectedComponentBase: string) => {
        setComponentBase(selectedComponentBase);
    };

    let skeletonBasesKeys: string | any[] = [];
    if (componentBases[selectedComponentBaseKey] !== undefined) {
        skeletonBasesKeys = Object.keys(componentBases[selectedComponentBaseKey]);
    }

    let selectedSkeletonBaseKey = '';
    if (skeletonBasesKeys.length > 0) {
        selectedSkeletonBaseKey = skeletonBasesKeys[0];
    }
    const [skeletonBaseName, setSkeletonBaseName] = React.useState(selectedSkeletonBaseKey);
    const onAccessSkeletonBase = (selectedSkeletonBaseName: string) => {
        setSkeletonBaseName(selectedSkeletonBaseName);
    };

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
                    <Box mt={2}>
                        <SelectedComponentBaseName componentBaseKeys={componentBaseKeys} componentBase={componentBase} onAccessComponentBase={onAccessComponentBase} />
                    </Box>
                    { cluster.componentBases[selectedComponentBaseKey] && skeletonBasesKeys.length > 0 &&
                        <Box mt={2}>
                            <SelectedSkeletonBaseName skeletonBaseKeys={skeletonBasesKeys} skeletonBaseName={skeletonBaseName} onAccessSkeletonBase={onAccessSkeletonBase}/>
                        </Box>
                    }
                    { cluster.componentBases[selectedComponentBaseKey] && skeletonBasesKeys.length > 0 && selectedSkeletonBaseKey != '' &&
                        <Box mt={2}>
                            <DefineDockerParams />
                        </Box>
                    }
                </Container>
            </Card>
        </div>
    )
}

export function SelectedSkeletonBaseName(props: any) {
    const {skeletonBaseName, skeletonBaseKeys, onAccessSkeletonBase} = props;

    console.log(skeletonBaseName, 'SelectedSkeletonBaseName')
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
                {skeletonBaseKeys.map((key: any) => (
                    <MenuItem key={key} value={key}>
                        {key}
                    </MenuItem>))
                }
            </Select>
        </FormControl>
    );
}