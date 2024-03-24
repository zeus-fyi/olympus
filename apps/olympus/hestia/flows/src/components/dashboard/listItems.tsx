import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import {Link} from "react-router-dom";
import {ExpandLess, ExpandMore} from "@mui/icons-material";
import {Collapse, List, ListSubheader} from "@mui/material";
import ConstructionIcon from '@mui/icons-material/Construction';
import CreditCardIcon from '@mui/icons-material/CreditCard';
import SecurityIcon from '@mui/icons-material/Security';
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import VpnKeyIcon from '@mui/icons-material/VpnKey';
import GraphicEqIcon from '@mui/icons-material/GraphicEq';
import {
    setOpenAiPanel,
    setOpenAppsPanel,
    setOpenClustersPanel,
    setOpenComputePanel
} from "../../redux/menus/menus.reducer";
import AutoFixNormalIcon from '@mui/icons-material/AutoFixNormal';

function SwapCallsIcon() {
    return null;
}

export default function MainListItems() {
    const [openServices, setOpenServices] = React.useState(false);
    const openClusters = useSelector((state: RootState) => state.menus.openClustersPanel);
    const openApps = useSelector((state: RootState) => state.menus.openAppsPanel);
    const isInternal = useSelector((state: RootState) => state.sessionState.isInternal);
    const openAiPanel = useSelector((state: RootState) => state.menus.openAiPanel);
    const openCompute = useSelector((state: RootState) => state.menus.openComputePanel);
    const dispatch = useDispatch();
    const handleClickServices = () => {
        setOpenServices(!openServices);
    };
    const handleClickClusters = () => {
        dispatch(setOpenClustersPanel(!openClusters));
    };
    const handleClickApps = () => {
        dispatch(setOpenAppsPanel(!openApps));
    };
    const handleClickCompute = () => {
        dispatch(setOpenComputePanel(!openCompute));
    };
    const handleClickAi = () => {
        dispatch(setOpenAiPanel(!openAiPanel));
    };
    return (
        <List
            sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}
            component="nav"
            aria-labelledby="nested-list-subheader"
            subheader={
                <ListSubheader component="div" id="nested-list-subheader">
                    Info Flows
                </ListSubheader>
            }
        >
            <div>
            <ListItemButton component={Link} onClick={handleClickAi} to="/ai">
                <ListItemIcon>
                    <GraphicEqIcon />
                </ListItemIcon>
                <ListItemText primary="AI" />
                {openAiPanel ? <ExpandLess  onClick={handleClickAi}/> : <ExpandMore onClick={handleClickAi}/>}
            </ListItemButton>
                <Collapse in={openAiPanel} timeout="auto" unmountOnExit>
                    <List component="div" disablePadding>
                        <ListItemButton sx={{ pl: 4 }} component={Link} to="/ai/workflow/builder">
                            <ListItemIcon>
                                <ConstructionIcon />
                            </ListItemIcon>
                            <ListItemText primary="Builder" />
                        </ListItemButton>
                    </List>
                    <List component="div" disablePadding>
                        <ListItemButton sx={{ pl: 4 }} component={Link} to="/ai/workflow/wizard">
                            <ListItemIcon>
                                <AutoFixNormalIcon />
                            </ListItemIcon>
                            <ListItemText primary="Flow Wizard" />
                        </ListItemButton>
                    </List>
                </Collapse>
            </div>
            {/*<ListItemButton component={Link} onClick={handleClickApps}  to="/apps">*/}
            {/*    <ListItemIcon>*/}
            {/*        <AppsIcon />*/}
            {/*    </ListItemIcon>*/}
            {/*    <ListItemText primary="Apps" />*/}
            {/*    {openApps ? <ExpandLess onClick={handleClickApps} /> : <ExpandMore onClick={handleClickApps}/>}*/}
            {/*</ListItemButton>*/}
            {/*<Collapse in={openApps} timeout="auto" unmountOnExit>*/}
            {/*    <List component="div" disablePadding>*/}
            {/*        <ListItemButton sx={{ pl: 4 }} component={Link} to="/apps/builder">*/}
            {/*            <ListItemIcon>*/}
            {/*                <ConstructionIcon />*/}
            {/*            </ListItemIcon>*/}
            {/*            <ListItemText primary="Builder" />*/}
            {/*        </ListItemButton>*/}
            {/*    </List>*/}
            {/*</Collapse>*/}
            {/*<ListItemButton onClick={handleClickCompute} component={Link} to="/compute/search">*/}
            {/*    <ListItemIcon>*/}
            {/*        <ManageSearchIcon />*/}
            {/*    </ListItemIcon>*/}
            {/*    <ListItemText primary="Compute" />*/}
            {/*    {openCompute ? <ExpandLess onClick={handleClickCompute} /> : <ExpandMore onClick={handleClickCompute}/>}*/}
            {/*</ListItemButton>*/}
            {/*<Collapse in={openCompute} timeout="auto" unmountOnExit>*/}
            {/*    <List component="div" disablePadding>*/}
            {/*        <ListItemButton sx={{ pl: 4 }} component={Link} to="/compute/summary">*/}
            {/*            <ListItemIcon>*/}
            {/*                <DnsIcon />*/}
            {/*            </ListItemIcon>*/}
            {/*            <ListItemText primary="Provisioned" />*/}
            {/*        </ListItemButton>*/}
            {/*    </List>*/}
            {/*</Collapse>*/}
            {/*<ListItemButton onClick={handleClickClusters} component={Link} to="/clusters">*/}
            {/*    <ListItemIcon>*/}
            {/*        <CloudIcon />*/}
            {/*    </ListItemIcon>*/}
            {/*    <ListItemText primary="Clusters"/>*/}
            {/*    {openClusters ? <ExpandLess onClick={handleClickClusters} /> : <ExpandMore onClick={handleClickClusters}/>}*/}
            {/*</ListItemButton>*/}
            {/*<Collapse in={openClusters} timeout="auto" unmountOnExit>*/}
            {/*    <List component="div" disablePadding>*/}
            {/*        <ListItemButton sx={{ pl: 4 }} component={Link} to="/clusters/config">*/}
            {/*            <ListItemIcon>*/}
            {/*                <TuneIcon />*/}
            {/*            </ListItemIcon>*/}
            {/*            <ListItemText primary="Config" />*/}
            {/*        </ListItemButton>*/}
            {/*    </List>*/}
            {/*</Collapse>*/}
            {/*<ListItemButton component={Link} to="/loadbalancing/dashboard">*/}
            {/*    <ListItemIcon>*/}
            {/*        <SwapCallsIcon />*/}
            {/*    </ListItemIcon>*/}
            {/*    <ListItemText primary="Load Balancing" />*/}
            {/*</ListItemButton>*/}
            {/*<ListItemButton onClick={handleClickServices}>*/}
            {/*    <ListItemIcon>*/}
            {/*        <ViewListIcon />*/}
            {/*    </ListItemIcon>*/}
            {/*    <ListItemText primary="Services" />*/}
            {/*    {openServices ? <ExpandLess /> : <ExpandMore />}*/}
            {/*</ListItemButton>*/}
            {/*<Collapse in={openServices} timeout="auto" unmountOnExit>*/}
            {/*    <List component="div" disablePadding>*/}
            {/*        {isInternal && (*/}
            {/*                <div>*/}
            {/*                    <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/mev">*/}
            {/*                        <ListItemIcon>*/}
            {/*                            <LeaderboardIcon />*/}
            {/*                        </ListItemIcon>*/}
            {/*                        <ListItemText primary="MEV" />*/}
            {/*                    </ListItemButton>*/}
            {/*                    <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/chatgpt">*/}
            {/*                        <ListItemIcon>*/}
            {/*                            <ChatIcon />*/}
            {/*                        </ListItemIcon>*/}
            {/*                        <ListItemText primary="ChatGPT" />*/}
            {/*                    </ListItemButton>*/}
            {/*                </div>*/}
            {/*            )*/}
            {/*        }*/}
            {/*        <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/ethereum/aws">*/}
            {/*            <ListItemIcon>*/}
            {/*                <AutoFixHighIcon />*/}
            {/*            </ListItemIcon>*/}
            {/*            <ListItemText primary="AWS Wizard" />*/}
            {/*        </ListItemButton>*/}
            {/*        <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/ethereum/validators">*/}
            {/*            <ListItemIcon>*/}
            {/*                <MiscellaneousServicesIcon />*/}
            {/*            </ListItemIcon>*/}
            {/*            <ListItemText primary="Validators" />*/}
            {/*        </ListItemButton>*/}
            {/*    </List>*/}
            {/*</Collapse>*/}
            <ListItemButton component={Link} to="/loadbalancing/dashboard">
                <ListItemIcon>
                    <SwapCallsIcon />
                </ListItemIcon>
                <ListItemText primary="Load Balancing" />
            </ListItemButton>
            <ListItemButton component={Link} to="/billing">
                <ListItemIcon>
                    <CreditCardIcon />
                </ListItemIcon>
                <ListItemText primary="Billing" />
            </ListItemButton>
            <ListItemButton component={Link} to="/secrets">
                <ListItemIcon>
                    <VpnKeyIcon />
                </ListItemIcon>
                <ListItemText primary="Secrets" />
            </ListItemButton>
            <ListItemButton component={Link} to="/access">
                <ListItemIcon>
                    <SecurityIcon />
                </ListItemIcon>
                <ListItemText primary="Access" />
            </ListItemButton>
        </List>
    );
}
