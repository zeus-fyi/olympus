import {Navigate} from "react-router-dom";
import Dashboard from "../dashboard/Dashboard";
import React, {useEffect} from "react";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {accessApiGateway} from "../../gateway/access";
import {setSessionAuth} from "../../redux/auth/session.reducer";

export const HomeLayout = () => {
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
    return <Dashboard />
};
