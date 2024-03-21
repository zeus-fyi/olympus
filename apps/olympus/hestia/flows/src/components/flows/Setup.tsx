import {Card, CardActionArea, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";

export function SetupCard(props: any) {
    return (
        <Card sx={{ maxWidth: 320 }}>
            <CardActionArea>
                <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large', fontWeight: 'thin', marginRight: '15px', color: '#151C2F' }}>
                        Setup
                    </Typography>
                </CardContent>
            </CardActionArea>
        </Card>
    );
}
