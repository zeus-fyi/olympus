let apiUrl = '';

class ConfigService  {
    constructor() {
        apiUrl = process.env.REACT_APP_BACKEND_ENDPOINT || 'http://localhost:9002';
    }

    get apiUrl() {
        return apiUrl;
    }
}

export const configService = new ConfigService();