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

export const App = () => {
    return (
            <Provider store={store}>
                <BrowserRouter>
                    <Routes>
                            <Route path="/" element={<HomeLayout />} />
                            <Route path="/login" element={<Login />} />
                            <Route path="/signup" element={<SignUp />} />
                            <Route path="/verify/email/:id" element={<VerifyEmail />} />
                        <Route>
                            <Route path="/dashboard" element={<ProtectedLayout />}/>
                            <Route>
                                <Route path="clusters" element={<Clusters />} />
                                <Route path="clusters/apps" element={<AppsPage />} />
                                <Route path="clusters/app/:id" element={<AppPageWrapper />} />
                                <Route path="clusters/builder" element={<ClusterBuilderPage />} />
                                <Route path="clusters/:id" element={<ClustersPage />} />
                                <Route path="apps/microservice" element={<AppPageWrapper app={"microservice"} />} />
                                <Route path="apps/avax" element={<AppPageWrapper app={"avax"} />} />
                                <Route path="apps/eth" element={<AppPageWrapper app={"ethereumEphemeralBeacons"} />} />
                            </Route>
                            <Route>
                                <Route path="services/chatgpt" element={<ChatGPTPage />} />
                                <Route path="services/ethereum/validators" element={<ValidatorsServices />} />
                                <Route path="services/ethereum/aws" element={<AwsWizard />} />
                            </Route>
                            <Route path="billing" element={<Billing />} />
                            <Route path="access" element={<Access />} />
                        </Route>
                    </Routes>
                </BrowserRouter>
            </Provider>
        );
}

