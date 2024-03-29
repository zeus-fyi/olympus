import {Box, Card, CardContent, Container, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";
import TextField from "@mui/material/TextField";
import {setClusterName} from "../../../../redux/clusters/clusters.builder.reducer";
import {AddComponentBases} from "./AddComponentBases";
import {DefineClusterComponentBaseParams} from "./DefineComponentBases";
import MdTextField from "./MdTutorial";

export function ClusterConfigPage(props: any) {
    const {} = props;
    return (
        <div>
            <Stack direction="row" spacing={2} sx={{ mb: 4 }}>
            <Stack direction="column" spacing={2} sx={{ mb: 4 }}>
            <Card sx={{ minWidth: 500, maxWidth: 500 }}>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Define Cluster Bases
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Define Cluster Class & Workload Bases
                    </Typography>
                </CardContent>
                <Container maxWidth="xl" sx={{ mb: 4 }}>
                    <Box mt={2}>
                        <ClusterName />
                    </Box>
                    <Box mt={2}>
                        <AddComponentBases />
                    </Box>
                </Container>
            </Card>
            <Box display="flex" flexDirection="row" sx={{ mt: 4 }}>
                <DefineClusterComponentBaseParams />
            </Box>
            </Stack>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <MdTextField />
            </Container>
        </Stack>
</div>
    );
}

export function ClusterName() {
    const dispatch = useDispatch();
    const clusterName = useSelector((state: RootState) => state.clusterBuilder.cluster.clusterName);
    const onClusterNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newClusterName = event.target.value;
        dispatch(setClusterName(newClusterName));
    };
    return (
        <TextField
            fullWidth
            id="clusterName"
            label="Cluster Name"
            variant="outlined"
            value={clusterName}
            onChange={onClusterNameChange}
            sx={{ width: '100%' }}
        />
    );
}