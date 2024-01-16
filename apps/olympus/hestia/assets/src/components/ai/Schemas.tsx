import {
    CardActions,
    CardContent,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    Stack,
    TextareaAutosize,
    ToggleButton
} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";
import Box from "@mui/material/Box";
import {setSchema} from "../../redux/ai/ai.reducer";
import {useDispatch} from "react-redux";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";

export function Schemas(props: any) {
    const {schema, createOrUpdateSchema, requestStatusSchema, requestStatusSchemaError} = props;

    const dispatch = useDispatch();
    return (
        <div>
            <Box sx={{ ml: 0, mr: 2 }} >
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Schemas
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Add JSON schema definitions for JSON tasks.
                    </Typography>
                </CardContent>
            </Box>
            <Stack direction="row" >
                <Box sx={{ width: '50%', ml:2, mb: 0, mt: 0 }}>
                    <TextField
                        label={`Schema Name`}
                        variant="outlined"
                        value={schema.schemaName}
                        onChange={(event) => dispatch(setSchema(
                            { ...schema, schemaName: event.target.value }))}
                        fullWidth
                    />
                </Box>
                <Box sx={{ width: '50%', ml:2, mb: 0, mt: 0 }}>
                    <TextField
                        label={`Schema Group`}
                        variant="outlined"
                        value={schema.schemaGroup}
                        onChange={(event) => dispatch(setSchema(
                            { ...schema, schemaGroup: event.target.value }))}
                        fullWidth
                    />
                </Box>
            </Stack>
            <Stack direction="row" >
                <Box  sx={{ mb: 2,ml: 2, mr:2, mt: 2  }}>
                    <ToggleButton
                        value="check"
                        selected={schema.isObjArray}
                        onChange={(event) => dispatch(setSchema(
                            { ...schema, isObjArray: !schema.isObjArray }))}
                    >
                        {schema.isObjArray ? 'JSON Object Array' : 'JSON Object'}
                    </ToggleButton>
                </Box>
                <Box  sx={{ mt: 4, mb: 0,ml: 2, mr:2  }}>
                    <Typography variant="body2" color="text.secondary">
                        Sets response type as JSON object or an array of JSON objects.
                    </Typography>
                </Box>
            </Stack>
            <Stack direction="row" >
                <Box flexGrow={7} sx={{ mb: 2,ml: 2, mr:0  }}>
                    <TextField
                        fullWidth
                        id="field-name"
                        label="Field Name"
                        variant="outlined"
                        // value={evalMetric.evalMetricName}
                        // onChange={(e) => dispatch(setEvalMetric({
                        //     ...evalMetric, // Spread the existing action properties
                        //     evalMetricName: e.target.value // Update the actionName
                        // }))}
                    />
                </Box>
                <Box flexGrow={7} sx={{ mb: 2,ml: 2, mr:0  }}>
                    <FormControl fullWidth >
                        <InputLabel id="field-data-type">Data Type</InputLabel>
                        <Select
                            labelId="field-data-type-label"
                            id="field-data-type-label"
                            // value={evalMetric.evalMetricDataType}
                            // label="Eval Metric Type"
                            // fullWidth
                            // onChange={(e) => dispatch(setEvalMetric({
                            //     ...evalMetric, // Spread the existing action properties
                            //     evalMetricDataType: e.target.value // Update the actionName
                            // }))}
                        >
                            <MenuItem value="number">{'number'}</MenuItem>
                            <MenuItem value="string">{'string'}</MenuItem>
                            <MenuItem value="boolean">{'boolean'}</MenuItem>
                            <MenuItem value="array[boolean]">{'array[boolean]'}</MenuItem>
                            <MenuItem value="array[number]">{'array[number]'}</MenuItem>
                            <MenuItem value="array[string]">{'array[string]'}</MenuItem>
                        </Select>
                    </FormControl>
                </Box>
                <Box sx={{ mt: 1, mb: 0,ml: 2, mr:0  }}>
                    <Button fullWidth variant={"contained"} >Add</Button>
                </Box>
                <Box sx={{ mt: 1, mb: 0,ml: 2, mr:2  }}>
                    <Button fullWidth variant={"contained"} >Clear</Button>
                    {/*<Button fullWidth variant={"contained"} onClick={clearEvalMetricRow}>Clear</Button>*/}
                </Box>
            </Stack>
            <Box  sx={{ mt: 1, mb: 0,ml: 2, mr:2  }}>
                <Typography variant="body2" color="text.secondary">
                    Field description
                </Typography>
            </Box>
                <Box  sx={{ mb: 2, mt: 2, ml: 2, mr: 2 }}>
                    <TextareaAutosize
                        minRows={18}
                        // value={editAggregateTask.prompt}
                        // onChange={(event) => dispatch(setEditAggregateTask({ ...editAggregateTask, prompt: event.target.value }))}
                        style={{ resize: "both", width: "100%" }}
                    />
                </Box>
            <CardActions>
                <Box flexGrow={1} sx={{ ml: 0, mr: 0}}>
                    <Button fullWidth variant="contained" onClick={createOrUpdateSchema}>Create or Update Schema</Button>
                </Box>
            </CardActions>
            {requestStatusSchema != '' && (
                <Container sx={{ mt: 2}}>
                    <Typography variant="h6" color={requestStatusSchemaError}>
                        {requestStatusSchema}
                    </Typography>
                </Container>
            )}
        </div>
    )
}