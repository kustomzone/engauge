const inMemoryJWTManager = () => {
    let logoutEventName = 'ra-logout';
    let refreshEndpoint = '/refresh-token';
    let inMemoryJWT = null;
    let refreshTimeOutId;

    window.addEventListener('storage', (event) => {
        if (event.key === logoutEventName) {
            inMemoryJWT = null;
        }
    });

    const setRefreshTokenEndpoint = endpoint => refreshEndpoint = endpoint;

    const refreshToken = (delay) => {
        refreshTimeOutId = window.setTimeout(
            getRefreshedToken,
            delay * 1000 - 5000
        ); // period in seconds, minus 5 seconds
    };

    const abortRefreshToken = () => {
        if (refreshTimeOutId) {
            window.clearTimeout(refreshTimeOutId);
        }
    };

    const getRefreshedToken = () => {
        const request = new Request(refreshEndpoint, {
            method: 'GET',
            headers: new Headers({ 'Content-Type': 'application/json' }),
            credentials: 'include',
        });
        return fetch(request)
            .then((response) => {
                if (response.status !== 200) {
                    eraseToken();
                    global.console.log(
                        'Failed to renew the jwt from the refresh token.'
                    );
                    return { token: null };
                }
                return response.json();
            })
            .then(({ token, tokenExpiry }) => {
                if (token) {
                    setToken(token, tokenExpiry);
                    return true;
                }

                return false;
            });
    };

    const getToken = () => inMemoryJWT;

    const setToken = (token, delay) => {
        inMemoryJWT = token;
        refreshToken(delay);
        return true;
    };

    const eraseToken = () => {
        inMemoryJWT = null;
        abortRefreshToken();
        window.localStorage.setItem(logoutEventName, Date.now());
        return true;
    }

    const setLogoutEventName = name => logoutEventName = name;

    return {
        eraseToken,
        getToken,
        setLogoutEventName,
        getRefreshedToken,
        setRefreshTokenEndpoint,
        setToken,
    }
};

export default inMemoryJWTManager();