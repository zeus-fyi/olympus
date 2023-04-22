class ConfigService  {
    private readonly apiUrl: string;
    private readonly zeusApiUrl: string;
    private readonly artemisApiUrl: string;
    private readonly heraApiUrl: string;
    private readonly stripePubKey: string;

    constructor() {
        this.apiUrl = process.env.REACT_APP_BACKEND_ENDPOINT || 'http://localhost:9002';
        this.zeusApiUrl = process.env.REACT_APP_ZEUS_BACKEND_ENDPOINT || 'http://localhost:9001';
        this.artemisApiUrl = process.env.REACT_APP_ARTEMIS_BACKEND_ENDPOINT || 'http://localhost:9004';
        this.heraApiUrl = process.env.REACT_APP_HERA_BACKEND_ENDPOINT || 'http://localhost:9008';
        this.stripePubKey = process.env.REACT_APP_STRIPE_PUBLISHABLE_KEY || 'pk_test_51MoIbzLLP9P61KzQDIlpiWOfoKF8CHJuHkLWjd01lQGfK8NrqCIUS9qS49j44g5AGK7J3g6064H4INbPn11zhsba00Bezb2Fop';
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
    public getHeraApiUrl(): string {
        return this.heraApiUrl;
    }
    public getStripePubKey(): string {
        return this.stripePubKey;
    }
}

export const configService = new ConfigService();