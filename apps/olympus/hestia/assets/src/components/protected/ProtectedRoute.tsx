import {Navigate} from "react-router-dom";
import Dashboard from "../dashboard/Dashboard";
import React from "react";

export const ProtectedRoute = () => {
    let parsedUser: any = localStorage.getItem('user');
    const user = JSON.parse(parsedUser);

    if (!user) {
        return <Navigate to="/login" />;
    }
    return <Dashboard />
};