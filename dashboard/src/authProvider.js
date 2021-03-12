import inMemoryJWTManager from "./inMemoryJwt";

const apiUrl = '';

const authProvider = {
    login: ({ username, password }) => {
        const request = new Request(`${apiUrl}/login`, {
            method: 'POST',
            body: JSON.stringify({ username, password }),
            headers: new Headers({ 'Content-Type': 'application/json' })
        });
        return fetch(request)
            .then((response) => {
                if (response.status < 200 || response.status >= 300) {
                    throw new Error(response.statusText);
                }

                return response.json();
            })
            .then(({ token }) => inMemoryJWTManager.setToken(token));
    },

    logout: () => {
        const request = new Request(`${apiUrl}/logout`, {
            method: 'GET',
            headers: new Headers({ 'Content-Type': 'application/json' }),
            credentials: 'include',
        });
        inMemoryJWTManager.eraseToken();

        return fetch(request).then(() => '/login');
    },

    checkAuth: () => {
        console.log('checkAuth');
        if (!inMemoryJWTManager.getToken()) {
            inMemoryJWTManager.setRefreshTokenEndpoint(`${apiUrl}/refresh-token`);
            return inMemoryJWTManager.getRefreshedToken().then(tokenHasBeenRefreshed => {
                return tokenHasBeenRefreshed ? Promise.resolve() : Promise.reject();
            });
        } else {
            return Promise.resolve();
        }
    },

    checkError: (error) => {
        const status = error.status;
        if (status === 401 || status === 403) {
            inMemoryJWTManager.eraseToken();
            return Promise.reject();
        }
        return Promise.resolve();
    },

    getPermissions: () => {
        return inMemoryJWTManager.getToken() ? Promise.resolve() : Promise.reject();
    },
};

export default authProvider;