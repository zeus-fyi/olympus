import {Card, CardActions, CardContent} from "@mui/material";
import * as React from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";

export function ZeusServiceRequest() {
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Create Zeus Validators Service Request
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates Zeus Validators Service Request
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Send</Button>
            </CardActions>
        </Card>
    );
}
