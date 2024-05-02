import {Card, CardContent, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";
import Checkbox from "@mui/material/Checkbox";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import TextField from "@mui/material/TextField";
import {setPreviewCount} from "../../redux/flows/flows.reducer";
import {useDispatch, useSelector} from "react-redux";

export function SetupCard(props: any) {
    const { webChecked, handleChangeWebChecked, vesChecked, handleChangeVesChecked, multiPromptOn, handleChangeMultiPromptOn, checked, checkedLi, handleChangeLi, handleChange, gs, handleChangeGs} = props;
    const previewCount = useSelector((state: any) => state.flows.previewCount);

    const handleChangeCount = (event: React.ChangeEvent<HTMLInputElement>) => {
        dispatch(setPreviewCount(Number(event.target.value)));
    };
    const dispatch = useDispatch()
    return (
        <div>
            <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 2}}>
            <div>
                <Card sx={{maxWidth: 320}}>
                    <CardContent style={{display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}>
                        <Typography gutterBottom variant="h5" component="div" style={{
                            fontSize: 'large',
                            fontWeight: 'thin',
                            marginRight: '15px',
                            color: '#151C2F'
                        }}>
                            Stages
                        </Typography>
                    </CardContent>
                    <Box sx={{mb: 2}}>
                        <Divider/>
                    </Box>
                    {/*<Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 0}}>*/}
                    {/*    <Checkbox*/}
                    {/*        checked={gs}*/}
                    {/*        onChange={handleChangeGs}*/}
                    {/*    />*/}
                    {/*    <Typography variant="body1">Google Search</Typography>*/}
                    {/*</Stack>*/}
                    <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 0, mr: 4}}>
                        <Checkbox
                            checked={checked}
                            onChange={handleChange}
                        />
                        <Typography variant="body1">LinkedIn Personal</Typography>
                    </Stack>
                    <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 0}}>
                        <Checkbox
                            checked={checkedLi}
                            onChange={handleChangeLi}
                        />
                        <Typography variant="body1">LinkedIn Biz</Typography>
                    </Stack>
                    <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 0}}>
                        <Checkbox
                            checked={vesChecked}
                            onChange={handleChangeVesChecked}
                        />
                        <Typography variant="body1">Validate Emails</Typography>
                    </Stack>
                    <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 2}}>
                        <Checkbox
                            checked={webChecked}
                            onChange={handleChangeWebChecked}
                        />
                        <Typography variant="body1">Fetch Website</Typography>
                    </Stack>
                    <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 2}}>
                        <TextField
                            label="Preview Count"
                            variant="outlined"
                            type="number"
                            value={previewCount}
                            inputProps={{ min: 1 }}  // Set minimum value to 0
                            onChange={handleChangeCount}
                        />
                        <Typography variant="body1">Preview Count</Typography>
                    </Stack>
                </Card>
            </div>
            {/*<Box sx={{mb: 2}}>*/}
            {/*    <Card sx={{maxWidth: 320}}>*/}
            {/*        <CardContent style={{display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}>*/}
            {/*            <Typography gutterBottom variant="h5" component="div" style={{*/}
            {/*                fontSize: 'large',*/}
            {/*                fontWeight: 'thin',*/}
            {/*                marginRight: '15px',*/}
            {/*                color: '#151C2F'*/}
            {/*            }}>*/}
            {/*                Config Options*/}
            {/*            </Typography>*/}
            {/*        </CardContent>*/}
            {/*        <Box sx={{mb: 2}}>*/}
            {/*            <Divider/>*/}
            {/*        </Box>*/}
            {/*        <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 0}}>*/}
            {/*            <Typography variant="body1">Multi-Prompt</Typography>*/}
            {/*            <Checkbox*/}
            {/*                checked={multiPromptOn}*/}
            {/*                onChange={handleChangeMultiPromptOn}*/}
            {/*            />*/}
            {/*        </Stack>*/}
            {/*    </Card>*/}
            {/*</Box>*/}
            </Stack>
        </div>
    );
}
