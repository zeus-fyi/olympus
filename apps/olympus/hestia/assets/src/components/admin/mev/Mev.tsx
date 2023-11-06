import * as React from 'react';
import {useEffect, useState} from 'react';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import {AppBar, Drawer} from '../../dashboard/Dashboard';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Container from '@mui/material/Container';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import Button from "@mui/material/Button";
import {useNavigate} from "react-router-dom";
import {useDispatch} from "react-redux";
import authProvider from "../../../redux/auth/auth.actions";
import MainListItems from "../../dashboard/listItems";
import {MevBundlesTable} from "./MevBundlesTable";
import {Card, CardContent, FormControl, InputLabel, MenuItem, Select, Stack, Tab, Tabs} from "@mui/material";
import {mevApiGateway} from "../../../gateway/mev";
import {MevCallBundlesTable} from "./MevCallBundlesTable";

const mdTheme = createTheme();

function MevContent(props: any) {
    const {callBundles,bundles, groups, groupName, selectedMainTab, handleMainTabChange} = props;
    const [open, setOpen] = React.useState(true);
    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();
    const dispatch = useDispatch();

    const handleLogout = (event: any) => {
        event.preventDefault();
        authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }

    return (
        <ThemeProvider theme={mdTheme}>
            <Box sx={{display: 'flex'}}>
                <CssBaseline/>
                <AppBar position="absolute" open={open} style={{backgroundColor: '#333'}}>
                    <Toolbar
                        sx={{
                            pr: '24px', // keep right padding when drawer closed
                        }}
                    >
                        <IconButton
                            edge="start"
                            color="inherit"
                            aria-label="open drawer"
                            onClick={toggleDrawer}
                            sx={{
                                marginRight: '36px',
                                ...(open && {display: 'none'}),
                            }}
                        >
                            <MenuIcon/>
                        </IconButton>
                        <Typography
                            component="h1"
                            variant="h6"
                            color="inherit"
                            noWrap
                            sx={{flexGrow: 1}}
                        >
                            MEV Dashboard
                        </Typography>
                        <Button
                            color="inherit"
                            onClick={handleLogout}
                        >Logout
                        </Button>
                    </Toolbar>
                </AppBar>
                <Drawer variant="permanent" open={open}>
                    <Toolbar
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'flex-end',
                            px: [1],
                        }}
                    >
                        <IconButton onClick={toggleDrawer}>
                            <ChevronLeftIcon/>
                        </IconButton>
                    </Toolbar>
                    <Divider/>
                    <List component="nav">
                        <MainListItems/>
                        <Divider sx={{my: 1}}/>
                    </List>
                </Drawer>
                <Box
                    component="main"
                    sx={{
                        backgroundColor: (theme) =>
                            theme.palette.mode === 'light'
                                ? theme.palette.grey[100]
                                : theme.palette.grey[900],
                        flexGrow: 1,
                        height: '100vh',
                        overflow: 'auto',
                    }}
                >
                    <Toolbar/>
                    <Container maxWidth={"xl"} sx={{mt: 4, mb: 4}}>
                        <Card className="onboarding-card-highlight-qn-routing-table"  sx={{ maxWidth: 700 }}>
                            <CardContent>
                                <Typography gutterBottom variant="h5" component="div">
                                   MEV Analytics
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                   Switch between bundles, tx analysis, and more.
                                </Typography>
                            </CardContent>
                            <Box mr={2} ml={2} mt={2} mb={4}>
                                <Stack direction={"row"} spacing={2} alignItems={"center"}>
                                    <FormControl sx={{  }} fullWidth variant="outlined">
                                        <InputLabel key={`groupNameLabel`} id={`groupName`}>
                                            Table View
                                        </InputLabel>
                                        <Select
                                            labelId={`groupNameLabel`}
                                            id={`groupName`}
                                            name="groupName"
                                            value={"bundles"}
                                            //onChange={(event) => handleChangeGroup(event.target.value)}
                                            label="Mev Group"
                                        >
                                            <MenuItem key={'all'} value={'-all'}>{"all"}</MenuItem>
                                            {Object.keys(groups).map((name) => <MenuItem key={name} value={name}>{name}</MenuItem>)}
                                        </Select>
                                    </FormControl>
                                </Stack>
                            </Box>
                        </Card>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        {groupName === "bundles" && (
                            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                                <Tabs value={selectedMainTab} onChange={handleMainTabChange} aria-label="basic tabs">
                                    <Tab className="onboarding-card-highlight-all-routes" label="Executed"  />
                                    <Tab className="onboarding-card-highlight-all-procedures" label="Simulated" />
                                </Tabs>
                            </Box>
                        )}
                    </Container>
                    { selectedMainTab == 0 && bundles && bundles.length > 0 &&
                    <Container maxWidth="xl" sx={{mt: 4, mb: 4}}>
                        <MevBundlesTable bundles={bundles}/>
                    </Container>
                    }
                    { selectedMainTab == 1 && callBundles && callBundles.length > 0 &&
                        <Container maxWidth="xl" sx={{mt: 4, mb: 4}}>
                            <MevCallBundlesTable callBundles={callBundles}/>
                        </Container>
                    }
                </Box>
            </Box>
        </ThemeProvider>
    );
}

type TraderInfoType = { [key: string]: { totalTxFees: number } };

export function createBundleData(
    eventID: string,
    submissionTime: string,
    bundleHash: string,
    bundledTxs: any[] = [],
    traderInfo: TraderInfoType,
    revenue: number,
    totalCost: number,
    totalGasCost: number,
    profit: number,
) {
    return {eventID, submissionTime, bundleHash, bundledTxs,traderInfo, revenue, totalCost, totalGasCost, profit};
}
export function createCallBundleData(
    eventID: string,
    submissionTime: string,
    bundleHash: string,
    builderName: string,
    bundleGasPrice: string,
    coinbaseDiff: string,
    gasFees: string,
    results: any[] = [],

) {
    return {eventID, submissionTime, bundleHash, builderName, bundleGasPrice, coinbaseDiff, gasFees, results};
}

export default function Mev() {
    const [bundles, setBundles] = useState([{}]);
    const [callBundles, setCallBundles] = useState([{}]);
    const [groupName, setGroupName] = useState('bundles');
    const [groups, setGroups] = useState({});
    const [loading, setLoading] = useState(true);
    const [selectedMainTab, setSelectedMainTab] = useState(0);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await mevApiGateway.getDashboardInfo();
                const mevDashboardTable: any[] = response.data.bundles;
                const mevDashboardTableRows = mevDashboardTable.map((v: any) =>
                    createBundleData(v.eventID, v.submissionTime, v.bundleHash, v.bundledTxs, v.traderInfo, v.revenue, v.totalCost, v.totalGasCost, v.profit)
                );
                setBundles(mevDashboardTableRows)

                const callBundlesTable: any[] = response.data.callBundles;
                const callBundlesTableRows = callBundlesTable.map((v: any) =>
                    createCallBundleData(v.eventID, v.submissionTime, v.flashbotsCallBundleResponse.bundleHash, v.builderName, v.flashbotsCallBundleResponse.bundleGasPrice, v.flashbotsCallBundleResponse.coinbaseDiff, v.flashbotsCallBundleResponse.gasFees, v.results)
                );
                console.log('sdsdf', callBundlesTableRows)
                setCallBundles(callBundlesTableRows)
                const mevTopKTokens: any[] = response.data.topKTokens;
                setGroups({
                    'bundles': bundles,
                    'callBundles': callBundles,
                    'topKTokens': mevTopKTokens
                })
            } catch (error) {
                console.log("error", error);
            } finally {
                setLoading(false);
            }
        }
        fetchData();
    }, []);
    if (loading) {
        return <div>Loading...</div>;
    }
    const handleMainTabChange = (event: React.SyntheticEvent, newValue: number) => {
        setSelectedMainTab(newValue);
    };
    return <MevContent callBundles={callBundles} bundles={bundles} groups={groups} groupName={groupName} selectedMainTab={selectedMainTab} handleMainTabChange={handleMainTabChange}/>;
}
