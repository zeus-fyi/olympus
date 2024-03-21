import * as React from "react";
import {Card, CardActionArea} from "@mui/material";
import {MbTaskCmdPrompt} from "./CommandPrompt";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import {SetupCard} from "./Setup";

export function Commands(props: any) {
    const [code, setCode] = React.useState("");
    return (
        <div>
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'left', justifyContent: 'center', mb: 2 }}>
                <SetupCard />
            </Box>
            <Card sx={{ maxWidth: 1200, justifyContent: 'center' }}>
            <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large',fontWeight: 'thin', marginRight: '15x', color: '#151C2F'}}>
                Agent Tasking Commands
            </Typography>
                <CardActionArea>
                    <MbTaskCmdPrompt language={"plaintext"} code={code} onChange={setCode} height={"200px"} width={"1200px"}/>
                    <Box mt={2} sx={{ display: 'flex', flexDirection: 'column', alignItems: 'right' }}>
                        <Button
                            variant="contained"
                            // onClick={onClickSubmit}
                            // disabled={buttonDisabledCreate}
                            sx={{ backgroundColor: '#00C48C', '&:hover': { backgroundColor: '#00A678' }}}
                        >
                            Submit
                        </Button>
                    </Box>
                </CardActionArea>
            </Card>

        </div>
    );
}