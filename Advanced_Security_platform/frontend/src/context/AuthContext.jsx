import { createContext, useContext, useState, useEffect } from 'react';
import api from '../lib/axios';
import { parseJwt } from '../lib/jwt';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    const checkAuth = async () => {
        const token = localStorage.getItem('access_token');
        if (!token) {
            setLoading(false);
            return;
        }

        try {
            // Decode token to get user info or call a /me endpoint if available.
            // Since we don't have a specific /me endpoint, we rely on the token validity
            // or we can decode the JWT to get role/email.
            // For now, let's assume we can basic decode or the dashboard fetch will fail.

            // A robust way: try to refresh to verify validity if needed,
            // or just assume valid until 401.

            // Let's decode the payload base64
            // Decode using the utility
            const payload = parseJwt(token);
            if (payload) {
                setUser({
                    email: payload.email,
                    role: payload.role,
                    id: payload.user_id
                });
            } else {
                throw new Error("Invalid token");
            }

        } catch (e) {
            console.error("Auth check failed", e);
            localStorage.removeItem('access_token');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        checkAuth();
    }, []);

    const login = async (email, password, captchaToken, otp = null) => {
        const payload = { email, password, captcha_token: captchaToken };
        if (otp) {
            payload.otp = otp;
        }
        const res = await api.post('/login', payload);
        return res.data; // Should return { access_token, refresh_token } or 202 message
    };

    const logout = async () => {
        try {
            await api.post('/user/logout');
        } catch (e) { /* ignore */ }
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        setUser(null);
        window.location.href = '/login';
    };

    return (
        <AuthContext.Provider value={{ user, setUser, login, logout, loading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
