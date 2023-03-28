import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";
import {AddSkeletonBases} from "./AddSkeletonBases";

export function DefineClusterComponentBaseParams(props: any) {
    const {} = props;
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);

    let defaultState = ''
    const componentBaseKeys = Object.keys(cluster.componentBases);
    if (cluster.componentBases !== null) {
        defaultState = componentBaseKeys[0]
    }
    const [componentBase, setComponentBase] = React.useState('');
    const onAccessComponentBase = (selectedComponentBase: string) => {
        setComponentBase(selectedComponentBase);
    };
    return (
        <div>
            <Card sx={{ maxWidth: 500 }}>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Define Cluster Component Base Elements
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Sets Cluster Component Base Elements
                    </Typography>
                </CardContent>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Box mt={2}>
                        <SelectedComponentBaseName componentBaseKeys={componentBaseKeys} componentBase={componentBase} onAccessComponentBase={onAccessComponentBase} />
                    </Box>
                    <Box mt={2}>
                        <AddSkeletonBases componentBaseName={componentBase}/>
                    </Box>
                </Container>
            </Card>
        </div>
    );
}

export function SelectedComponentBaseName(props: any) {
    const {componentBase, componentBaseKeys, onAccessComponentBase} = props;

    return (
        <FormControl variant="outlined" style={{ minWidth: '100%' }}>
            <InputLabel id="network-label">Component Bases</InputLabel>
            <Select
                labelId="componentBase-label"
                id="componentBase"
                value={componentBase}
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