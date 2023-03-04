import axios from 'axios';
import {configService} from "../../config/config";

export default axios.create({
    baseURL: configService.apiUrl,
});