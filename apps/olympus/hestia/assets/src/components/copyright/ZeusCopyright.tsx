import Typography from "@mui/material/Typography";
import Link from "@mui/material/Link";
import * as React from "react";

export function ZeusCopyright(props: any) {
    return (
        <Typography variant="body2" color="text.secondary" align="center" {...props}>
            {'Copyright © Zeusfyi, Inc'}
            <Link color="inherit" href="https://cloud.zeus.fyi">
                Zeusfyi Cloud
            </Link>{' '}
            {new Date().getFullYear()}
            {'.'}
        </Typography>
    );
}