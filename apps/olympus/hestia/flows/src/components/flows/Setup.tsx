import {Card, CardContent, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";
import Checkbox from "@mui/material/Checkbox";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";

export function SetupCard(props: any) {
    const { checked, handleChange, gs, handleChangeGs} = props;
    return (
        <Card sx={{ maxWidth: 320 }}>
                <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large', fontWeight: 'thin', marginRight: '15px', color: '#151C2F' }}>
                        Stages
                    </Typography>
                </CardContent>
            <Box sx={{mb: 2}}>

             <Divider  />
            </Box>
            <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2}}>
                <Typography variant="body1">Google Search</Typography>
                <Box sx={{ml: 2, mb: 2}}>
                </Box>
                <Checkbox
                    checked={gs}
                    onChange={handleChangeGs}
                />
            </Stack>
            <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 2}}>
                <Typography variant="body1">LinkedIn</Typography>
                <Checkbox
                    checked={checked}
                    onChange={handleChange}
                />
            </Stack>
        </Card>
    );
}
