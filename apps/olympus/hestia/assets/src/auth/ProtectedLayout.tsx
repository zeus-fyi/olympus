import {Navigate} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {accessApiGateway} from "../gateway/access";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../redux/store";
import {setInternalAuth, setIsBillingSetup, setSessionAuth} from "../redux/auth/session.reducer";
import {setUserPlanDetails} from "../redux/loadbalancing/loadbalancing.reducer";

export const ProtectedLayout = (props: any) => {
    const {children} = props;
    const sessionAuthed = useSelector((state: RootState) => state.sessionState.sessionAuth);
    const [loading, setLoading] = useState(true);

    const dispatch = useDispatch();
    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await accessApiGateway.checkAuth();
                if (response.status >= 300 || response.status < 200) {
                    dispatch(setSessionAuth(false));
                    return;
                }
                if (response.data.planUsageDetails != null){
                    dispatch(setUserPlanDetails(response.data.planUsageDetails))
                }
                if (response.data.isBillingSetup === true) {
                    dispatch(setIsBillingSetup(true));
                }
                dispatch(setSessionAuth(true));
                if (response.data.isInternal === true) {
                    dispatch(setInternalAuth(true));
                } else {
                    dispatch(setInternalAuth(false));
                }
            } catch (error) {
                dispatch(setSessionAuth(false));
                setLoading(false);
            } finally {
                setLoading(false);
            }
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

export const InternalProtectedLayout = (props: any) => {
    const {children} = props;
    const sessionAuthed = useSelector((state: RootState) => state.sessionState.sessionAuth);
    const [loading, setLoading] = useState(true);

    const dispatch = useDispatch();
    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await accessApiGateway.checkAuth();
                if (response.status >= 300 || response.status < 200) {
                    dispatch(setSessionAuth(false));
                    return;
                }
                if (response.data.planUsageDetails != null){
                    dispatch(setUserPlanDetails(response.data.planUsageDetails))
                }
                if (response.data.isBillingSetup === true) {
                    dispatch(setIsBillingSetup(true));
                }
                dispatch(setSessionAuth(true));
                if (response.data.isInternal === true) {
                    dispatch(setInternalAuth(true));
                } else {
                    dispatch(setInternalAuth(false));
                }
            } catch (error) {
                dispatch(setSessionAuth(false));
                setLoading(false);
            } finally {
                setLoading(false);
            }
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
