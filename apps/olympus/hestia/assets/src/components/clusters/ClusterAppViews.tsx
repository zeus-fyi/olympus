import {CardContent, Stack, Switch} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";

export const ClusterViews = (props: any) => {
    const { pageView, setPageView } = props;
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPageView(event.target.checked);
    };

    return (
        <div>
            <Stack direction={"row"} spacing={2} alignItems={"center"}>
                {pageView ? (
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Cluster Apps View
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            This view shows all the apps that are deployed
                        </Typography>
                    </CardContent>
                ) : (
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Cluster View
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            This view shows all the cloud clusters you have and can interact with.
                        </Typography>
                    </CardContent>
                )}
                <Switch
                    checked={pageView}
                    onChange={handleChange}
                    color="primary"
                    name="pageView"
                    inputProps={{ 'aria-label': 'toggle page view' }}
                />
            </Stack>
        </div>
    );
};