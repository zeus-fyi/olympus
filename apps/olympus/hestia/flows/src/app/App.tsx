import React from 'react';
import {Provider} from 'react-redux';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import './App.css';
import store from "../redux/store";
import Login from "../components/login/Login";
import {InternalProtectedLayout, ProtectedLayout} from "../auth/ProtectedLayout";
import {HomeLayout} from "../components/home/Home";
import SignUp from "../components/signup/Signup";
import {VerifyEmail} from "../components/signup/VerifyEmail";
import Billing from "../components/billing/Billing";
import Access from "../components/access/Access";
import {VerifyQuickNodeLoginJWT} from "../components/login/VerifyLoginJWT";
import {GoogleOAuthProvider} from '@react-oauth/google';
import {configService} from "../config/config";
import ReactGA from "react-ga4";
import AiWorkflowsDashboard from "../components/ai/AI";
import Secrets from "../components/secrets/Secrets";
import AiWorkflowsEngineBuilderDashboard from "../components/ai/WorkflowBuilder";
import {TwitterLogin} from "../components/login/TwitterLogin";
import BizAutomationWizard from "../components/flows/Wizard";
import LoadBalancingDashboard from "../components/loadbalancing/LoadBalancingDashboard";

export const App = () => {
    ReactGA.initialize([
        {
            trackingId: "G-KZFWQL2CJN",
            // gaOptions: { 'debug_mode':true }, // optional
            //gtagOptions: {...}, // optional
        },
    ]);
    return (
        <GoogleOAuthProvider clientId={configService.getGoogClientID()}>
        <Provider store={store}>
                <BrowserRouter>
                    <Routes>
                            <Route path="/" element={<HomeLayout />} />
                            <Route path="/social/v1/twitter/callback" element={<TwitterLogin />} />
                            <Route path="/login" element={<Login />} />
                            <Route path="/signup" element={<SignUp />} />
                        <Route path="/quicknode/dashboard" element={<VerifyQuickNodeLoginJWT />} />
                        <Route path="/quicknode/access" element={<VerifyQuickNodeLoginJWT />} />
                        <Route path="/verify/email/:id" element={<VerifyEmail />} />
                        <Route>
                            <Route path="ai" element={<ProtectedLayout children={<AiWorkflowsDashboard />}/>}/>
                            <Route path="ai/workflow/builder" element={<ProtectedLayout children={<AiWorkflowsEngineBuilderDashboard />}/>}/>
                            <Route path="ai/workflow/wizard" element={<ProtectedLayout children={<BizAutomationWizard />}/>}/>
                            <Route path="ai/admin" element={<InternalProtectedLayout children={<BizAutomationWizard />}/>}/>
                            <Route path="billing" element={<ProtectedLayout children={<Billing />}/>}/>
                            <Route path="access" element={<ProtectedLayout children={<Access />}/>}/>
                            <Route path="secrets" element={<ProtectedLayout children={<Secrets />}/>}/>
                        </Route>
                        <Route path="loadbalancing/dashboard" element={<ProtectedLayout children={<LoadBalancingDashboard />}/>}/>
                    </Routes>
                </BrowserRouter>
            </Provider>
        </GoogleOAuthProvider>);
}

