import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import * as React from "react";

// TODO: Implement export function from results

export function CsvExportButton(props: any) {
    const { onUpload } = props;
    return (
        <Stack direction="row" alignItems="center" sx={{mt: 2}}>
            <Button variant="contained" component="label" style={{ backgroundColor: '#8991B0', color: '#151C2F' }}>
                <CloudUploadIcon />
                <input
                    hidden
                    accept="text/csv, application/json"
                    type="file"
                    onChange={onUpload}
                />
            </Button>
        </Stack>
    );
}