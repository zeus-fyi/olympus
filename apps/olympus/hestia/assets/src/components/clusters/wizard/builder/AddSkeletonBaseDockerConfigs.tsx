import {useDispatch, useSelector} from "react-redux";
import * as React from "react";
import {Card, CardContent, Container} from "@mui/material";
import {RootState} from "../../../../redux/store";
import {SelectedComponentBaseName} from "./DefineComponentBases";
import Typography from "@mui/material/Typography";

export function AddSkeletonBaseDockerConfigs(props: any) {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const componentBases = cluster.componentBases;
    const componentBaseKeys = Object.keys(componentBases);


    let selectedKey = '';
    if (componentBaseKeys.length > 0) {
        selectedKey = componentBaseKeys[0];
    }
    const [componentBase, setComponentBase] = React.useState(selectedKey);
    const onAccessComponentBase = (selectedComponentBase: string) => {
        setComponentBase(selectedComponentBase);
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
                    <SelectedComponentBaseName componentBaseKeys={componentBaseKeys} componentBase={componentBase} onAccessComponentBase={onAccessComponentBase} />
                </Container>
            </Card>
        </div>
    )
}
