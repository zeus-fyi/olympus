import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";

export function DefineClusterComponentBaseParams(props: any) {
    const {} = props;
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
                        <SelectedComponentBaseName />
                    </Box>
                </Container>
            </Card>
        </div>
    );
}

export function SelectedComponentBaseName(props: any) {
    const dispatch = useDispatch();
    const [componentBase, setComponentBase] = React.useState('');

    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const componentBaseKeys = Object.keys(cluster.componentBases);

    const onAccessComponentBase = (selectedComponentBase: string) => {
        setComponentBase(selectedComponentBase);
    };
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
                {componentBaseKeys.map((key) => (
                    <MenuItem key={key} value={key}>
                        {key}
                    </MenuItem>))
                }
            </Select>
        </FormControl>
    );
}