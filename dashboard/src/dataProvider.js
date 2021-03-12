import { fetchUtils } from 'react-admin';
import { stringify } from 'query-string';
import inMemoryJWTManager from "./inMemoryJwt";

const apiUrl = '';

const getClient = (url) => {
    const options = {
        headers: new Headers({ Accept: 'application/json' }),
    };
    const token = inMemoryJWTManager.getToken();

    if (token) {
        options.headers.set('Authorization', `Bearer ${token}`);
        return fetchUtils.fetchJson(url, options);
    } else {
    }
};

const postClient = (resource, params) => {
    const url = `${apiUrl}/dashboard/${resource}`
    const options = {
        method: 'POST',
        body: JSON.stringify(params.data),
        headers: new Headers({ Accept: 'application/json' }),
    };
    const token = inMemoryJWTManager.getToken();

    if (token) {
        options.headers.set('Authorization', `Bearer ${token}`);
        return fetchUtils.fetchJson(url, options);
    } else {
    }
};

const putClient = (resource, params) => {
    var url = `${apiUrl}/dashboard/${resource}/${params.id}`
    if (resource === "settings") {
        url = `${apiUrl}/dashboard/${resource}`
    }
    
    const options = {
        method: 'PUT',
        body: JSON.stringify(params.data),
        headers: new Headers({ Accept: 'application/json' }),
    };
    const token = inMemoryJWTManager.getToken();

    if (token) {
        options.headers.set('Authorization', `Bearer ${token}`);
        return fetchUtils.fetchJson(url, options);
    } else {
    }
};

const deleteClient = (resource, params) => {
    const url = `${apiUrl}/dashboard/${resource}/${params.id}`
    const options = {
        method: 'DELETE',
        headers: new Headers({ Accept: 'application/json' }),
    };
    const token = inMemoryJWTManager.getToken();

    if (token) {
        options.headers.set('Authorization', `Bearer ${token}`);
        return fetchUtils.fetchJson(url, options);
    } else {
    }
};

export default {
    getList: (resource, params) => {
        const { page, perPage } = params.pagination;
        const offset = (page - 1) * perPage
        const limit = perPage
        const url = `${apiUrl}/dashboard/${resource}?offset=${offset}&limit=${limit}`;

        return getClient(url, {
            method: 'GET',
        }).then(({ headers, json }) => ({
            data: json,
            total: parseInt(
                headers.get('x-total-count')
            ),
        }));
    },

    getOne: (resource, params) => {
        if (resource === "settings") {
            return getClient(`${apiUrl}/dashboard/${resource}`, {
                method: 'GET',
            }).then(({ json }) => ({
                data: json,
            }));
        };

        return getClient(`${apiUrl}/dashboard/${resource}/${params.id}`, {
            method: 'GET',
        }).then(({ json }) => ({
            data: json,
        }));
    },

    getMany: (resource, params) => {
        const url = `${apiUrl}/dashboard/${resource}`;
        return getClient(url).then(({ json }) => ({ data: json }));
    },

    getManyReference: (resource, params) => {
        const url = `${apiUrl}/dashboard/${resource}`;
        return getClient(url).then(({ headers, json }) => ({
            data: json,
            total: parseInt(
                headers.get('x-total-count'),
                10
            ),
        }));
    },

    update: (resource, params) =>
        putClient(resource, params).then(({ json }) => ({ data: json })),

    updateMany: (resource, params) => {
        const query = {
            filter: JSON.stringify({ id: params.ids }),
        };
        return getClient(`${apiUrl}/dashboard/${resource}?${stringify(query)}`, {
            method: 'PUT',
            body: JSON.stringify(params.data),
        }).then(({ json }) => ({ data: json }));
    },

    create: (resource, params) =>
        postClient(resource, params).then(({ json }) => ({ data: json })),

    delete: (resource, params) =>
        deleteClient(resource, params).then(({ json }) => ({ data: json })),

    deleteMany: (resource, params) => {
        const query = {
            filter: JSON.stringify({ id: params.ids }),
        };
        return getClient(`${apiUrl}/dashboard/${resource}?${stringify(query)}`, {
            method: 'DELETE',
            body: JSON.stringify(params.data),
        }).then(({ json }) => ({ data: json }));
    }
};
