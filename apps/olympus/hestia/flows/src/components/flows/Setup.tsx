import {Card, CardContent, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";
import Checkbox from "@mui/material/Checkbox";

export function SetupCard(props: any) {
    const { checked, handleChange, gs, handleChangeGs} = props;
    return (
        <Card sx={{ maxWidth: 320 }}>
                <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large', fontWeight: 'thin', marginRight: '15px', color: '#151C2F' }}>
                        Stages
                    </Typography>
                </CardContent>
            <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2}}>
                <Typography variant="body1">LinkedIn</Typography>
                <Checkbox
                    checked={checked}
                    onChange={handleChange}
                />
            </Stack>
            <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2}}>
                <Typography variant="body1">Google Search</Typography>
                <Checkbox
                    checked={gs}
                    onChange={handleChangeGs}
                />
            </Stack>
        </Card>
    );
}
