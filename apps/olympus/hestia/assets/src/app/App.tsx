import React, {useEffect} from 'react';
import {Provider} from 'react-redux';
import {BrowserRouter, Route, Routes, useLocation} from 'react-router-dom';
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

declare global {
    interface Window {
        dataLayer: any[];
        gtag: (...args: any[]) => void;
    }
}

export const App = () => {
    const location = useLocation();
    useEffect(() => {
        // Send virtual pageviews to Google Analytics on location change
        window.gtag('config', 'G-KZFWQL2CJN', {
            'page_path': location.pathname + location.search,
        });
    }, [location]);
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
                                <Route path="clusters/:id" element={<ClustersPage />} />
                            </Route>
                            <Route>
                                <Route path="services/ethereum/validators" element={<ValidatorsServices />} />
                                <Route path="services/ethereum/aws" element={<AwsWizard />} />
                            </Route>
                            {/*<Route path="access" element={<Access />} />*/}
                            {/*<Route path="billing" element={<Billing />} />*/}
                        </Route>
                    </Routes>
                </BrowserRouter>
            </Provider>
        );
}

