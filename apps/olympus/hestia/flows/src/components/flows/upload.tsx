import {Card, CardActionArea, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import * as React from "react";

export function CsvUploadActionAreaCard(props: any) {
    const onUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
        console.log('onUpload', 'event', event.target.files);
        const files = event.target.files;
        const file = files?.item(0);
        if (!file) return;

        const reader = new FileReader();
        reader.onload = (e) => {
            const result = e.target?.result;
            if (!result) {
                console.error("File read resulted in null.");
                return;
            }

            let data;
            if (file.type === "application/json") {
                try {
                    data = JSON.parse(result as string);
                } catch (error) {
                    console.error("Error parsing JSON file:", error);
                    return;
                }
            } else if (file.type === "text/csv") {
                try {
                    data = parseCSV(result as string);
                } catch (error) {
                    console.error("Error parsing CSV file:", error);
                    return;
                }
            } else {
                console.error("Unsupported file type:", file.type);
                return;
            }
        };

        reader.readAsText(file);
    };
    return (
        <Card sx={{ maxWidth: 320 }}>
            <CardActionArea>
                <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}>
                    <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large',fontWeight: 'thin', marginRight: '15x', color: '#151C2F'}}>
                        Upload CSV
                    </Typography>
                    <UploadButton onUpload={onUpload}/>
                </CardContent>
            </CardActionArea>
        </Card>
    );
}
export function UploadButton(props: any) {
    const { onUpload } = props;

    return (
        <Stack direction="row" alignItems="center" spacing={2}>
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




const parseCSV = (csvText: string): any[] => {
    const lines = csvText.split(/\r\n|\n/);
    const headers = lines[0].split(',');
    return lines.slice(1).map((line: string) => {
        const data = line.split(',');
        return headers.reduce((obj: { [key: string]: string }, nextKey: string, index: number) => {
            obj[nextKey] = data[index];
            return obj;
        }, {});
    });
};