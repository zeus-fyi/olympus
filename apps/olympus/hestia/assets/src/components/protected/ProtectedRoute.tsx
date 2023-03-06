import {Navigate} from "react-router-dom";
import Dashboard from "../dashboard/Dashboard";
import React from "react";

export const ProtectedRoute = () => {
    let userID = localStorage.getItem('userID');
    if (!userID) {
        return <Navigate to="/login" />;
    }
    return <Dashboard />
};