import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react'
import axios from 'axios';
import {configService} from "../../config/config";
import {RootState} from "../../redux/store";

export default axios.create({
    baseURL: configService.apiUrl,
});

export const zeusApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl:  configService.zeusApiUrl,
        prepareHeaders: (headers, { getState }) => {
            // By default, if we have a token in the store, let's use that for authenticated requests
            let state = getState() as RootState;
            const token = state.auth.token
            if (token) {
                headers.set('authorization', `Bearer ${token}`)
            }
            return headers
        },
    }),
    endpoints(build) {
        return {
            getClusters: build.query({ query: () => ({ url: '/v1/infra/read/org/topologies', method: 'get' }) }),
        }
    },
})

