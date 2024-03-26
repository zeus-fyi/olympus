import {Card, CardContent, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";
import Checkbox from "@mui/material/Checkbox";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";

export function SetupCard(props: any) {
    const { webChecked, handleChangeWebChecked, vesChecked, handleChangeVesChecked, multiPromptOn, handleChangeMultiPromptOn, checked, checkedLi, handleChangeLi, handleChange, gs, handleChangeGs} = props;
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
                    <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 0}}>
                        <Typography variant="body1">Google Search</Typography>
                        <Checkbox
                            checked={gs}
                            onChange={handleChangeGs}
                        />
                    </Stack>
                    <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 0}}>
                        <Typography variant="body1">LinkedIn Personal</Typography>
                        <Checkbox
                            checked={checked}
                            onChange={handleChange}
                        />
                    </Stack>
                    {/*<Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 2}}>*/}
                    {/*    <Typography variant="body1">LinkedIn Business</Typography>*/}
                    {/*    <Checkbox*/}
                    {/*        checked={checkedLi}*/}
                    {/*        onChange={handleChangeLi}*/}
                    {/*    />*/}
                    {/*</Stack>*/}
                    <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 0}}>
                        <Typography variant="body1">Validate Emails</Typography>
                        <Checkbox
                            checked={vesChecked}
                            onChange={handleChangeVesChecked}
                        />
                    </Stack>
                    <Stack direction="row" alignItems="center" spacing={2} sx={{ml: 2, mb: 2}}>
                        <Typography variant="body1">Fetch Website</Typography>
                        <Checkbox
                            checked={webChecked}
                            onChange={handleChangeWebChecked}
                        />
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
