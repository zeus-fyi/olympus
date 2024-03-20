import * as React from 'react';
import Button from '@mui/material/Button';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import Stack from '@mui/material/Stack';
import {Card, CardActionArea, CardContent, CardMedia} from "@mui/material";
import Typography from "@mui/material/Typography";

export function ValidatorsUploadActionAreaCard(props: any) {
    const { onValidatorsDepositsUpload } = props;

    return (
        <Card sx={{ maxWidth: 320 }}>
            <CardActionArea>
                <CardMedia
                    component="img"
                    height="230"
                    image={require("../../static/ethereum-logo.png")}
                    alt="ethereum"
                />
                <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', backgroundColor: '#8991B0'}}>
                    <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large',fontWeight: 'thin', marginRight: '15x', color: '#151C2F'}}>
                        Upload Deposits
                    </Typography>
                    <UploadValidatorsButton onValidatorsDepositsUpload={onValidatorsDepositsUpload}/>
                </CardContent>
            </CardActionArea>
        </Card>
    );
}
export function UploadValidatorsButton(props: any) {
    const { onValidatorsDepositsUpload } = props;

    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <Button variant="contained" component="label" style={{ backgroundColor: '#8991B0', color: '#151C2F' }}>
                <CloudUploadIcon />
                <input
                    hidden
                    accept="application/json"
                    type="file"
                    onChange={onValidatorsDepositsUpload}
                />
            </Button>
        </Stack>
    );
}