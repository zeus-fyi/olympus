import Typography from "@mui/material/Typography";
import * as React from "react";

export function ZeusCopyright(props: any) {
    return (
        <Typography variant="body2" color="text.secondary" align="center" {...props}>
            {'\nZeusfyi, Inc \n'}
            {' © '}
            {new Date().getFullYear()}
        </Typography>
    );
}