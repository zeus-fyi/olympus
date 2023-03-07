import React from 'react';
import {Provider} from 'react-redux';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import './App.css';
import store from "../redux/store";
import Login from "../components/login/Login";
import {ProtectedLayout} from "../auth/ProtectedLayout";
import {HomeLayout} from "../components/home/Home";
import Integrations from "../components/integrations/Integrations";
import ValidatorsServices from "../components/validators/Validators";
import Clusters from "../components/clusters/Clusters";
import ClustersPage from "../components/clusters/ClusterPage";

export const App = () => {
    return (
            <Provider store={store}>
                <BrowserRouter>
                    <Routes>
                            <Route path="/" element={<HomeLayout />} />
                            <Route path="/login" element={<Login />} />
                        <Route>
                            <Route path="/dashboard" element={<ProtectedLayout />}/>
                            <Route>
                                <Route path="clusters" element={<Clusters />} />
                                <Route path="clusters/:id" element={<ClustersPage />} />
                            </Route>
                            <Route>
                                <Route path="services/ethereum/validators" element={<ValidatorsServices />} />
                            </Route>
                            <Route path="integrations" element={<Integrations />} />
                        </Route>
                    </Routes>
                </BrowserRouter>
            </Provider>
        );
}

