import axios from 'axios';
import {configService} from "../../config/config";

export const awsLambdaInvoke = axios.create({
    baseURL: '',
});

export const hestiaApi = axios.create({
    baseURL: configService.apiUrl,
});

export const artemisApi = axios.create({
    baseURL: configService.artemisApiUrl,
});

export const zeusApi = axios.create({
    baseURL: configService.zeusApiUrl,
});