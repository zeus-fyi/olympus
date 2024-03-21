import {Card, CardActionArea, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import * as React from "react";
import {useDispatch} from "react-redux";
import {setPromptHeaders, setUploadTasksContent} from "../../redux/flows/flows.reducer";
import Container from "@mui/material/Container";
import {TaskPromptsTable} from "./PromptTable";
import {PromptsTextFieldRows} from "./UploadFieldMap";

export function AnalyzeActionAreaCard(props: any) {
    const dispatch = useDispatch();
    const onUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
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
            let headers;
            if (file.type === "application/json") {
                try {
                    data = JSON.parse(result as string);
                    // Assuming you want to set headers for JSON as well,
                    // you might need to derive them from the data structure
                    // For example, if data is an array of objects:
                    // headers = Object.keys(data[0]);
                } catch (error) {
                    console.error("Error parsing JSON file:", error);
                    return;
                }
            } else if (file.type === "text/csv") {
                try {
                    // Correctly destructure data and headers from the parseCSV result
                    const parseResult = parseCSV(result as string);
                    data = parseResult.data;
                    headers = parseResult.fields;
                    // console.log(data, headers);
                    dispatch(setPromptHeaders(headers));
                } catch (error) {
                    console.error("Error parsing CSV file:", error);
                    return;
                }
            } else {
                console.error("Unsupported file type:", file.type);
                return;
            }
            dispatch(setUploadTasksContent(data));
        };
        reader.readAsText(file);
    };
    return (
        <div>
            <Card sx={{ maxWidth: 400 }}>
                <CardActionArea>
                    <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}>
                        <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large',fontWeight: 'thin', marginRight: '15x', color: '#151C2F'}}>
                            Upload Prompts CSV
                        </Typography>
                        <UploadButton onUpload={onUpload}/>
                    </CardContent>
                </CardActionArea>
            </Card>
            <Container maxWidth="xl" sx={{ ml: -5, mt: 4}}>
                <PromptsTextFieldRows/>
            </Container>
            <Container maxWidth="xl" sx={{ ml: -5, mt: 4}}>
                <TaskPromptsTable/>
            </Container>
        </div>
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

type CsvDataRow = {
    [key: string]: string;
};
const parseCSV = (csvText: string): { data: CsvDataRow[], fields: string[] } => {
    const lines = csvText.split(/\r\n|\n/);
    // Ensure there's at least one line for headers
    if (lines.length === 0) {
        return { data: [], fields: [] };
    }

    const headers = lines[0].split(',');
    const data = lines.slice(1).filter(line => line).map(line => {
        const values = line.split(',');
        // Use the CsvDataRow type for the object
        const rowData: CsvDataRow = {};
        headers.forEach((header, index) => {
            rowData[header] = values[index] || ''; // Assign empty string if value is undefined
        });
        return rowData;
    });

    return { data, fields: headers };
};
