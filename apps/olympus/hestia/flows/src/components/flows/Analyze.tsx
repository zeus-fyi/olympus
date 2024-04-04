import {Card, CardActionArea, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import * as React from "react";
import {useDispatch} from "react-redux";
import {setPromptHeaders, setPromptsCsvContent, setUploadContacts} from "../../redux/flows/flows.reducer";
import Container from "@mui/material/Container";
import {TaskPromptsTable} from "./PromptTable";
import {PromptsTextFieldRows} from "./UploadFieldMap";
import Papa from "papaparse";

export function AnalyzeActionAreaCard(props: any) {
    const dispatch = useDispatch();
    const onUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
        const files = event.target.files;
        const file = files?.item(0);
        if (!file) return;

        if (file.type === "application/json") {
            const reader = new FileReader();
            reader.onload = (e) => {
                const result = e.target?.result;
                if (!result) {
                    console.error("File read resulted in null.");
                    return;
                }

                try {
                    const data = JSON.parse(result as string);
                    // Assuming you want to set headers for JSON as well,
                    // you might need to derive them from the data structure
                    // For example, if data is an array of objects:
                    // headers = Object.keys(data[0]);
                    dispatch(setUploadContacts(data));
                    // You should dispatch the headers here if necessary
                } catch (error) {
                    console.error("Error parsing JSON file:", error);
                    return;
                }
            };
            reader.readAsText(file);
        } else if (file.type === "text/csv") {
            Papa.parse(file, {
                complete: (result) => {
                    try {
                        const data = result.data;
                        const headers = result.meta.fields || [];
                       if (Array.isArray(headers)) {
                           dispatch(setPromptsCsvContent(data as []))
                           dispatch(setPromptHeaders(headers))
                            // dispatch(setCsvHeaders(headers));
                            // dispatch(setUploadContacts(data as []));
                        }
                    } catch (error) {
                        console.error("Error parsing CSV file:", error);
                    }
                },
                header: true
            });
        } else {
            console.error("Unsupported file type:", file.type);
        }
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
            <Button variant="contained" component="label" style={{  backgroundColor: '#4fd3ad', color: '#FFF' }}>
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

