import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";
import {useState} from "react";
import {AddSkeletonBases} from "./AddSkeletonBases";
import {
    setSelectedComponentBaseName,
    setSelectedContainerName,
    setSelectedSkeletonBaseName
} from "../../../../redux/clusters/clusters.builder.reducer";

export function DefineClusterComponentBaseParams(props: any) {
    const {} = props;
    const [viewField, setViewField] = useState('');
    const onChangeComponentOrSkeletonBase = () => {
        setViewField('')
    }
    return (
        <div>
            <Card sx={{ minWidth: 500, maxWidth: 500 }}>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Define Cluster Workloads
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        You'll configure these workloads next. Give them a name for now.
                    </Typography>
                </CardContent>
                <Container maxWidth="xl" sx={{ mb: 4 }}>
                    <Box mt={2}>
                        <SelectedComponentBaseName onChangeComponentOrSkeletonBase={onChangeComponentOrSkeletonBase}/>
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
    const {onChangeComponentOrSkeletonBase} = props;
    const dispatch = useDispatch();
    let cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    let selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const onAccessComponentBase = (selectedComponentBaseName: string) => {
       dispatch(setSelectedComponentBaseName(selectedComponentBaseName));
       const skeletonBaseName = Object.keys(cluster.componentBases[selectedComponentBaseName])[0];
       dispatch(setSelectedSkeletonBaseName(skeletonBaseName));
        // Add a check to see if the `containers` field exists
        if (cluster.componentBases[selectedComponentBaseName] &&
            cluster.componentBases[selectedComponentBaseName][skeletonBaseName] &&
            cluster.componentBases[selectedComponentBaseName][skeletonBaseName].containers) {
            const containerKeys = Object.keys(cluster.componentBases[selectedComponentBaseName][skeletonBaseName].containers);
            if (containerKeys.length > 0) {
                dispatch(setSelectedContainerName(containerKeys[0]));
            }
        }
        onChangeComponentOrSkeletonBase();
    };

    let show = Object.keys(cluster.componentBases).length > 0;
    return (
        <div>
            {show &&
            <FormControl sx={{mb: 1}} variant="outlined" style={{ minWidth: '100%' }}>
                <InputLabel id="network-label">Cluster Bases</InputLabel>
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
