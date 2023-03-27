import {Box, Card, CardContent, Container} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";
import TextField from "@mui/material/TextField";
import {AddPortsInputFields} from "./DefinePorts";

export function DefineDockerParams(props: any) {
    const {} = props;
    return (
        <div>
            <Card sx={{ maxWidth: 500 }}>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Set Docker Image
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Sets Docker Image Default
                    </Typography>
                </CardContent>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Box mt={2}>
                        <DockerImageName />
                    </Box>
                    <Box mt={2}>
                        <AddPortsInputFields />
                    </Box>
                </Container>
            </Card>
        </div>
    );
}

export function DockerImageName() {
    const dispatch = useDispatch();
    const dockerImageName = useSelector((state: RootState) => state.clusterBuilder.cluster.clusterName);
    const onDockerImageNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newDockerImageName = event.target.value;
        //dispatch(addDockerImagePort(newDockerImageName));
    };
    return (
        <TextField
            fullWidth
            id="dockerImage"
            label="Docker Image Name"
            variant="outlined"
            value={dockerImageName}
            onChange={onDockerImageNameChange}
            sx={{ width: '100%' }}
        />
    );
}