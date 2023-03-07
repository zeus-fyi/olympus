let apiUrl = '';
let zeusApiUrl = '';

class ConfigService  {
    constructor() {
        apiUrl = process.env.REACT_APP_BACKEND_ENDPOINT || 'http://localhost:9002';
        zeusApiUrl = process.env.REACT_APP_ZEUS_BACKEND_ENDPOINT || 'http://localhost:9001';
    }
    get apiUrl() {
        return apiUrl;
    }
    get zeusApiUrl() {
        return zeusApiUrl;
    }
}

export const configService = new ConfigService();