import type React from 'react';
import { createContext, useState, useEffect, useCallback, useMemo } from 'react';
import { apiGet, apiPost } from '../utils/apiClient';
import { logger } from '../utils/logger';

export interface User {
    authenticated: boolean;
    username: string;
    isSuperuser?: boolean;
    expiresAt?: string;
}

export interface AuthContextValue {
    user: User | null;
    loading: boolean;
    login: (username: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
    forceLogout: () => Promise<void>;
}

interface AuthProviderProps { children: React.ReactNode }

interface UserInfoResponse {
    authenticated: boolean;
    username: string;
    is_superuser?: boolean;
}

interface LoginResponse {
    success: boolean;
    message?: string;
    expires_at?: string;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export const AuthProvider = ({ children }: AuthProviderProps): React.ReactElement => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => { void checkAuth(); }, []);

    const checkAuth = async (): Promise<void> => {
        try {
            const info = await apiGet<UserInfoResponse>('/api/v1/user/info');
            if (info.authenticated) {
                setUser({
                    authenticated: true,
                    username: info.username,
                    isSuperuser: info.is_superuser ?? false,
                });
            } else {
                setUser(null);
            }
        } catch (error) {
            logger.error('Auth check failed:', error);
            setUser(null);
        } finally {
            setLoading(false);
        }
    };

    const login = useCallback(async (username: string, password: string): Promise<void> => {
        const result = await apiPost<LoginResponse>('/api/v1/auth/login', { username, password });
        if (!result.success) {
            throw new Error(result.message ?? 'Authentication failed');
        }
        try {
            const info = await apiGet<UserInfoResponse>('/api/v1/user/info');
            setUser({
                authenticated: true,
                username: info.username,
                isSuperuser: info.is_superuser ?? false,
                expiresAt: result.expires_at,
            });
        } catch {
            setUser({
                authenticated: true,
                username,
                expiresAt: result.expires_at,
            });
        }
    }, []);

    const logout = useCallback(async (): Promise<void> => {
        try {
            await apiPost('/api/v1/auth/logout');
        } catch (error) {
            logger.error('Logout request failed:', error);
        }
        setUser(null);
    }, []);

    const forceLogout = useCallback(async (): Promise<void> => {
        try { await apiPost('/api/v1/auth/logout'); } catch { /* ignore */ }
        setUser(null);
    }, []);

    const value: AuthContextValue = useMemo(() => ({
        user, loading, login, logout, forceLogout,
    }), [user, loading, login, logout, forceLogout]);

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export default AuthContext;
