import React from 'react';
import {Provider} from 'react-redux';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import './App.css';
import store from "../redux/store";
import Login from "../components/login/Login";
import {ProtectedRoute} from "../components/protected/ProtectedRoute";
import {HomeLayout} from "../components/home/Home";

export const App = () => {
    return (
            <Provider store={store}>
                <BrowserRouter>
                    <Routes>
                            <Route path="/" element={<HomeLayout />} />
                            <Route path="/login" element={<Login />} />
                        <Route>
                            <Route path="/dashboard" element={<ProtectedRoute />}/>
                        </Route>
                    </Routes>
                </BrowserRouter>
            </Provider>
        );
}

