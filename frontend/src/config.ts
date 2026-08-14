// Backend base URL, configurable at build time via VITE_BACKEND_URL (e.g. "https://remi-backend.fly.dev").
// Falls back to the local dev backend (same host, port 8080) when unset.
const configuredBackendUrl = import.meta.env.VITE_BACKEND_URL?.replace(/\/+$/, '');

export function getApiBaseUrl(): string {
    if (configuredBackendUrl) return configuredBackendUrl;
    return `${window.location.protocol}//${window.location.hostname}:8080`;
}

export function getWsBaseUrl(): string {
    return getApiBaseUrl().replace(/^http/, 'ws');
}
