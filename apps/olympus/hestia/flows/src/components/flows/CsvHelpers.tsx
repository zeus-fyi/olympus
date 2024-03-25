import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import CloudDownloadIcon from "@mui/icons-material/CloudUpload";
import * as React from "react";
import {useState} from "react";

// TODO: Implement export function from results
const filterAggregationTypes = (data: any) => {
    if (!data || !data.aggregatedData) {
        // Return an empty array or handle the error as appropriate
        return [];
    }

    return data.aggregatedData
        .filter((item: any) => item.taskType === "aggregation")
        .map((item: any) => item.completionChoices);
};

const CsvExportButton = (props: any) => {
    const { name, results } = props;
    const [data, setData] =  useState(filterAggregationTypes(results));
    console.log('data:', data)
    return (
        <Stack direction="row" alignItems="center" spacing={2} sx={{mt: 2}}>
            <Button
                variant="contained"
                style={{ backgroundColor: '#4CAF50', color: '#FFF' }}
                onClick={() => parseJSONAndCreateCSV(name, data)}
            >
                <CloudDownloadIcon />
            </Button>
        </Stack>
    );
};

export default CsvExportButton;

export const parseJSONAndCreateCSV = (name: string, data: any) => {
    const processedData = prettyPrintObject(data);

    const sja =  stringToJsonArray(processedData)
    const csvContent = jsonArrayToCSV(sja);
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    downloadBlobAsFile( name+'.csv', blob)
};

const jsonArrayToCSV = (jsonArray: any[]): string => {
    if (jsonArray.length === 0) return '';

    // Extract headers
    const headers = Object.keys(jsonArray[0])
        .map(key => key.charAt(0).toUpperCase() + key.slice(1)) // Capitalize the first letter
        .join(',');
    // Extract rows
    const rows = jsonArray.map(obj =>
        Object.values(obj).map((value: any) =>
            `"${value.toString().replace(/"/g, '""')}"` // Escape double quotes
        ).join(',')
    );

    return [headers, ...rows].join('\n');
};
function stringToJsonArray(jsonString: string) {
    // Assuming each JSON object is separated by a newline
    // Split the string by newline to get an array of strings, each representing a JSON object
    const jsonParts = jsonString.trim().split('\n');

    // Parse each part as JSON and return the array of objects
    return jsonParts.map(part => JSON.parse(part));
}

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

const convertToCSV = (arr: any[]) => {
    // Check if there's data to convert
    if (arr.length === 0) {
        return '';
    }

    // Extract headers
    const headers = Object.keys(arr[0]).join(',');
    // Map each object's values, ensuring to handle commas within values
    const rows = arr.map(obj =>
        Object.values(obj).map((value: any) =>
            `"${value.toString().replace(/"/g, '""')}"` // Escape double quotes
        ).join(',')
    );

    return [headers, ...rows].join('\n');
};
