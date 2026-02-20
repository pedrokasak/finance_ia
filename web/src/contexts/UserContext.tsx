import { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { api } from '@/api/client';

interface UserProfile {
    id: string;
    email: string;
    first_name?: string;
    last_name?: string;
    avatar_url?: string;
    plan?: string;
    financial_method_id?: string;
    notifications_enabled?: boolean;
    two_fa_enabled?: boolean;
}

interface UserContextType {
    profile: UserProfile | null;
    setProfile: (p: UserProfile | null) => void;
    refreshProfile: () => Promise<void>;
    updateAvatar: (base64: string) => void;
}

const UserContext = createContext<UserContextType | null>(null);

function getJWTPayload(): { email?: string; plan?: string; user_id?: string } {
    try {
        const token = localStorage.getItem('authToken');
        if (!token) return {};
        return JSON.parse(atob(token.split('.')[1]));
    } catch {
        return {};
    }
}

export function UserProvider({ children }: { children: ReactNode }) {
    const [profile, setProfile] = useState<UserProfile | null>(() => {
        const jwt = getJWTPayload();
        if (!jwt.user_id) return null;
        return { id: jwt.user_id, email: jwt.email || '', plan: jwt.plan || 'free' };
    });

    const refreshProfile = useCallback(async () => {
        const jwt = getJWTPayload();
        if (!jwt.user_id) return;
        try {
            const res = await api.get(`/user/${jwt.user_id}`);
            const u = (res as any).data?.user || (res as any).data;
            if (u) {
                setProfile({
                    id: u.id || jwt.user_id,
                    email: u.email || jwt.email || '',
                    first_name: u.first_name,
                    last_name: u.last_name,
                    avatar_url: u.avatar_url,
                    plan: u.plan || jwt.plan || 'free',
                    financial_method_id: u.financial_method_id,
                    notifications_enabled: u.notifications_enabled,
                    two_fa_enabled: u.two_fa_enabled,
                });
            }
        } catch {
            // Keep cached data
        }
    }, []);

    const updateAvatar = useCallback((base64: string) => {
        setProfile((prev) => prev ? { ...prev, avatar_url: base64 } : prev);
    }, []);

    return (
        <UserContext.Provider value={{ profile, setProfile, refreshProfile, updateAvatar }}>
            {children}
        </UserContext.Provider>
    );
}

export function useUser() {
    const ctx = useContext(UserContext);
    if (!ctx) throw new Error('useUser must be used inside UserProvider');
    return ctx;
}
