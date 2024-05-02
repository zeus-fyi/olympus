import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import CloudDownloadIcon from "@mui/icons-material/CloudUpload";
import * as React from "react";
import {aiApiGateway} from "../../gateway/ai";

const CsvExportButton = (props: any) => {
    const [loading, setIsLoading] = React.useState(false);
    const { name, orchStrID, results } = props;
    const onClickExportCsv = async (name: string, id: string) => {
        try {
            setIsLoading(true);
            const response = await aiApiGateway.flowCsvExportRequest(id);
            // Assuming response is the CSV string
            const blob = new Blob([response.data], { type: 'text/csv' });
            downloadBlobAsFile(`${name}`, blob);
        } finally {
            setIsLoading(false);
        }
    };

    if (loading) {
        return <div>Loading...</div>
    }
    return (
        <Stack direction="row" alignItems="center" spacing={2} sx={{mt: 2}}>
            <Button
                variant="contained"
                style={{ backgroundColor: '#4fd3ad', color: '#FFF' }}
                onClick={() => onClickExportCsv(name, orchStrID)}
            >
                <CloudDownloadIcon />
            </Button>
        </Stack>
    );
};

export default CsvExportButton;

export const prettyPrintObject = (obj: any): string => {
    // If 'message' is a string that needs to be parsed
    try {
        if (Array.isArray(obj)) {
            return obj.map(obj => prettyPrintObject(obj)).join('\n');
        } else if (typeof obj === 'string') {
            return prettyPrintObject(JSON.parse(obj))
        } else if (obj.prompt) {
            return prettyPrintObject(obj.prompt);
        } else if (obj.content) {
            return prettyPrintObject(obj.content);
        } else if (obj.tool_calls) {
            return prettyPrintObject(obj.tool_calls);
        } else if (obj.tool_uses) {
            return prettyPrintObject(obj.tool_uses);
        } else if (obj.arguments) {
            return prettyPrintObject(obj.arguments);
        } else if (obj.message) {
            return prettyPrintObject(obj.message);
        } else if (obj.function) {
            return prettyPrintObject(obj.function);
        } else if (obj.parameters) {
            return prettyPrintObject(obj.parameters);
        } else if (obj['google-search-results-agg']) {
            return prettyPrintObject(obj['google-search-results-agg']);
        } else if (obj['results-agg']) {
            return prettyPrintObject(obj['results-agg']);
        } else {
            return JSON.stringify(obj, null,0);
        }
    } catch (error) {
        // Return the original string if it can't be parsed
        return obj;
    }
};


export const downloadBlobAsFile = (fileName: string, blob: Blob) => {
    // Create a URL for the blob
    const url = window.URL.createObjectURL(blob);

    // Create an invisible <a> element with a link to the blob
    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = url;
    a.download = fileName; // Set the file name for the download

    // Append the <a> element to the body
    document.body.appendChild(a);

    // Simulate a click on the <a> element
    a.click();

    // Clean up by removing the <a> element and revoking the blob URL
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
};
