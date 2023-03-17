class ConfigService  {
    private readonly apiUrl: string;
    private readonly zeusApiUrl: string;
    private readonly artemisApiUrl: string;

    constructor() {
        this.apiUrl = process.env.REACT_APP_BACKEND_ENDPOINT || 'http://localhost:9002';
        this.zeusApiUrl = process.env.REACT_APP_ZEUS_BACKEND_ENDPOINT || 'http://localhost:9001';
        this.artemisApiUrl = process.env.REACT_APP_ARTEMIS_BACKEND_ENDPOINT || 'http://localhost:9004';
    }
    public getApiUrl(): string {
        return this.apiUrl;
    }

    public getZeusApiUrl(): string {
        return this.zeusApiUrl;
    }

    public getArtemisApiUrl(): string {
        return this.artemisApiUrl;
    }
}

export const configService = new ConfigService();