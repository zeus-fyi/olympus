import {Navigate, useOutlet} from "react-router-dom";
import React, {useEffect} from "react";
import Dashboard from "../components/dashboard/Dashboard";
import {accessApiGateway} from "../gateway/access";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../redux/store";
import {setSessionAuth} from "../redux/auth/session.reducer";

export const ProtectedLayout = () => {
    const outlet = useOutlet();
    const sessionAuthed = useSelector((state: RootState) => state.sessionState.sessionAuth);
    const dispatch = useDispatch();
    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await accessApiGateway.checkAuth();
                if (response.status !== 200) {
                    dispatch(setSessionAuth(false));
                    return;
                }
                dispatch(setSessionAuth(true));
            } catch (error) {
                dispatch(setSessionAuth(false));
                console.log("error", error);
            }}
        fetchData().then(r =>
            console.log("")
        );
    }, [sessionAuthed]);

    if (!sessionAuthed) {
        return <Navigate to="/login" />;
    }
    return (
        <div>
            <Dashboard />
            {outlet}
        </div>
    );
};