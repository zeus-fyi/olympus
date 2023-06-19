import {Navigate, useOutlet} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {accessApiGateway} from "../gateway/access";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../redux/store";
import {setSessionAuth} from "../redux/auth/session.reducer";

export const ProtectedLayout = (props: any) => {
    const outlet = useOutlet();
    const {children} = props;
    const sessionAuthed = useSelector((state: RootState) => state.sessionState.sessionAuth);
    const [loading, setLoading] = useState(true);

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
                setLoading(false);
            }
            setLoading(false);
        }
        fetchData().then(r =>
            console.log("")
        );
    }, []);
    if (loading) {
        return null;
    }
    if (!sessionAuthed) {
        return <Navigate to="/login" />;
    }
    return (
        <div>
            {children}
        </div>
    );
};