import {useDispatch, useSelector} from "react-redux";
import * as React from "react";
import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import {RootState} from "../../../../redux/store";
import {SelectedComponentBaseName} from "./DefineComponentBases";
import Typography from "@mui/material/Typography";

export function AddSkeletonBaseDockerConfigs(props: any) {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const componentBases = cluster.componentBases;
    const componentBaseKeys = Object.keys(componentBases);
    console.log(cluster)

    let selectedComponentBaseKey = '';
    if (componentBaseKeys.length > 0) {
        selectedComponentBaseKey = componentBaseKeys[0];
    }
    const [componentBase, setComponentBase] = React.useState(selectedComponentBaseKey);
    const onAccessComponentBase = (selectedComponentBase: string) => {
        setComponentBase(selectedComponentBase);
    };

    const skeletonBasesKeys = Object.keys(componentBases[selectedComponentBaseKey]);

    console.log(skeletonBasesKeys)
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
            <Card sx={{ maxWidth: 500 }}>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Define Cluster Skeleton Base Elements
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Sets Cluster Skeleton Base Elements
                    </Typography>
                </CardContent>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Box mt={2}>
                        <SelectedComponentBaseName componentBaseKeys={componentBaseKeys} componentBase={componentBase} onAccessComponentBase={onAccessComponentBase} />
                    </Box>
                    { cluster.componentBases[selectedComponentBaseKey] &&
                        <Box mt={2}>
                            <SelectedSkeletonBaseName skeletonBaseKeys={skeletonBasesKeys} skeletonBaseName={skeletonBaseName} onAccessSkeletonBase={onAccessSkeletonBase}/>
                        </Box>
                    }
                </Container>
            </Card>
        </div>
    )
}

export function SelectedSkeletonBaseName(props: any) {
    const {skeletonBaseName, skeletonBaseKeys, onAccessSkeletonBase} = props;

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