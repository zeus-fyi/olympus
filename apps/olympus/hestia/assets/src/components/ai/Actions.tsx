import {Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import * as React from "react";


export function Actions(props: any) {

    return (
        <div>
            <Typography gutterBottom variant="h5" component="div">
                Actions
            </Typography>
            <Typography variant="body2" color="text.secondary">
                Actions are the steps that will be performed on the data.
            </Typography>
            <Stack direction="row">
                <Box flexGrow={3} sx={{width: '50%', mb: 0, mt: 2, mr: 1}}>
                    <TextField
                        label={`Action Name`}
                        variant="outlined"
                        fullWidth
                    />
                </Box>
                <Box flexGrow={3} sx={{width: '50%', mb: 0, mt: 2, ml: 1}}>
                    <TextField
                        label={`Action Group`}
                        variant="outlined"
                        fullWidth
                    />
                </Box>
            </Stack>
            <Stack direction="row">
            </Stack>
        </div>
    )
}