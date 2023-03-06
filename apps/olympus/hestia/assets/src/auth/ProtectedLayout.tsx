import {Navigate, useOutlet} from "react-router-dom";
import React from "react";
import Dashboard from "../components/dashboard/Dashboard";

export const ProtectedLayout = () => {
    let userID = localStorage.getItem('userID');
    const outlet = useOutlet();

    if (!userID) {
        return <Navigate to="/login" />;
    }
    return (
        <div>
            <Dashboard />
            {outlet}
        </div>
    );
};