interface JWTManager {
    ereaseToken: () => boolean;
    getRefreshedToken: () => Promise<boolean>;
    getToken: () => string | null;
    setLogoutEventName: (name: string) => void;
    setRefreshTokenEndpoint: (endpoint: string) => void;
    setToken: (token: string, delay: number) => boolean;
    waitForTokenRefresh: () => Promise<boolean>;
}

const inMemoryJWTManager = (): JWTManager => {
    let inMemoryJWT: string | null = null;
    let isRefreshing: Promise<boolean | void>;
    let logoutEventName = 'ra-logout';
    let refreshEndpoint = '/refresh-token';
    let refreshTimeOutId: number | undefined;

    const setLogoutEventName = (name: string) => {
        logoutEventName = name;
    };

    const setRefreshTokenEndpoint = (endpoint: string) => {
        refreshEndpoint = endpoint;
    };

    // This countdown feature is used to renew the JWT before it's no longer valid
    // in a way that is transparent to the user.
    const refreshToken = (delay: number): void => {
        refreshTimeOutId = window.setTimeout(
            getRefreshedToken,
            delay * 1000 - 5000
        ); // Validity period of the token in seconds, minus 5 seconds
    };

    const abordRefreshToken = (): void => {
        if (refreshTimeOutId) {
            window.clearTimeout(refreshTimeOutId);
        }
    };

    const waitForTokenRefresh = (): Promise<boolean> => {
        if (!isRefreshing) {
            return Promise.resolve(false);
        }
        return isRefreshing.then(() => {
            return true;
        });
    };

    // The method make a call to the refresh-token endpoint
    // If there is a valid cookie, the endpoint will set a fresh jwt in memory.
    const getRefreshedToken = (): Promise<boolean> => {
        const request = new Request(refreshEndpoint, {
            method: 'GET',
            headers: new Headers({ 'Content-Type': 'application/json' }),
            credentials: 'include',
        });

        isRefreshing = fetch(request)
            .then((response) => {
                if (response.status !== 200) {
                    ereaseToken();
                    console.log('Token renewal failure');
                    return { token: null };
                }
                return response.json();
            })
            .then(({ token, tokenExpiry }) => {
                if (token) {
                    setToken(token, tokenExpiry);
                    return true;
                }
                ereaseToken();
                return false;
            });

        return isRefreshing.then(() => true);
    };

    const getToken = (): string | null => inMemoryJWT;

    const setToken = (token: string, delay: number): boolean => {
        inMemoryJWT = token;
        refreshToken(delay);
        return true;
    };

    const ereaseToken = (): boolean => {
        inMemoryJWT = null;
        abordRefreshToken();
        window.localStorage.setItem(logoutEventName, Date.now().toString());
        return true;
    };

    // This listener will allow to disconnect a session of ra started in another tab
    window.addEventListener('storage', (event) => {
        if (event.key === logoutEventName) {
            inMemoryJWT = null;
        }
    });

    return {
        ereaseToken,
        getRefreshedToken,
        getToken,
        setLogoutEventName,
        setRefreshTokenEndpoint,
        setToken,
        waitForTokenRefresh,
    };
};

let InMemoryJWT: JWTManager;
export default InMemoryJWT = inMemoryJWTManager();