import React from 'react';
import {Provider} from 'react-redux';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import './App.css';
import store from "../redux/store";
import Login from "../components/login/Login";
import {ProtectedLayout} from "../auth/ProtectedLayout";
import {HomeLayout} from "../components/home/Home";
import ValidatorsServices from "../components/validators/Validators";
import Clusters from "../components/clusters/Clusters";
import ClustersPage from "../components/clusters/ClusterPage";
import AwsWizard from "../components/validators/AwsWizard";
import SignUp from "../components/signup/Signup";
import {VerifyEmail} from "../components/signup/VerifyEmail";
import ClusterBuilderPage from "../components/clusters/wizard/ClusterBuilderPage";
import AppsPage from "../components/apps/AppsPage";
import {AppPageWrapper} from "../components/app/AppPageWrapper";
import Billing from "../components/billing/Billing";
import Access from "../components/access/Access";
import {ChatGPTPage} from "../components/chatgpt/ChatGPTWrapper";
import Dashboard from "../components/dashboard/Dashboard";
import {VerifyQuickNodeLoginJWT} from "../components/login/VerifyLoginJWT";

export const App = () => {
    return (
            <Provider store={store}>
                <BrowserRouter>
                    <Routes>
                            <Route path="/" element={<HomeLayout />} />
                            <Route path="/login" element={<Login />} />
                            <Route path="/signup" element={<SignUp />} />
                        <Route path="/quicknode/dashboard" element={<VerifyQuickNodeLoginJWT />} />
                        <Route path="/quicknode/access" element={<VerifyQuickNodeLoginJWT />} />
                        <Route path="/verify/email/:id" element={<VerifyEmail />} />
                        <Route>
                            <Route path="/dashboard" element={<ProtectedLayout children={<Dashboard />}/>}/>
                            <Route>
                                <Route path="clusters"  element={<ProtectedLayout children={<Clusters />}/>}/>
                                <Route path="clusters/apps" element={<ProtectedLayout children={<AppsPage />}/>}/>
                                <Route path="clusters/app/:id" element={<ProtectedLayout children={<AppPageWrapper />}/>}/>
                                <Route path="clusters/builder" element={<ProtectedLayout children={<ClusterBuilderPage />}/>}/>
                                <Route path="clusters/:id" element={<ProtectedLayout children={<ClustersPage />}/>}/>
                                <Route path="apps/microservice" element={<ProtectedLayout children={<AppPageWrapper app={"microservice"} />}/>}/>
                                <Route path="apps/avax" element={<ProtectedLayout children={<AppPageWrapper app={"avax"} />}/>} />
                                <Route path="apps/eth" element={<ProtectedLayout children={<AppPageWrapper app={"ethereumEphemeralBeacons"} />}/>} />
                            </Route>
                            <Route>
                                <Route path="services/quicknode/dashboard" element={<Dashboard />} />
                                <Route path="services/chatgpt" element={<ProtectedLayout children={<ChatGPTPage />}/>}/>
                                <Route path="services/ethereum/validators" element={<ProtectedLayout children={<ValidatorsServices />}/>}/>
                                <Route path="services/ethereum/aws" element={<ProtectedLayout children={<AwsWizard />}/>}/>
                            </Route>
                            <Route path="billing" element={<ProtectedLayout children={<Billing />}/>}/>
                            <Route path="access"  element={<ProtectedLayout children={<Access />}/>}/>
                        </Route>
                    </Routes>
                </BrowserRouter>
            </Provider>
        );
}

