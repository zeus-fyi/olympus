let apiUrl = '';
let zeusApiUrl = '';
let artemisApiUrl = '';

class ConfigService  {
    constructor() {
        apiUrl = process.env.REACT_APP_BACKEND_ENDPOINT || 'http://localhost:9002';
        zeusApiUrl = process.env.REACT_APP_ZEUS_BACKEND_ENDPOINT || 'http://localhost:9001';
        artemisApiUrl = process.env.REACT_APP_ARTEMIS_BACKEND_ENDPOINT || 'http://localhost:9004';
    }
    get apiUrl() {
        return apiUrl;
    }
    get zeusApiUrl() {
        return zeusApiUrl;
    }
    get artemisApiUrl() {
        return artemisApiUrl;
    }
}

export const configService = new ConfigService();