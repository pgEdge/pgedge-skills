import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { AuthProvider } from '../AuthContext';
import { useAuth } from '../useAuth';

const wrapper = ({ children }: { children: React.ReactNode }) => (
    <AuthProvider>{children}</AuthProvider>
);

describe('AuthContext', () => {
    beforeEach(() => {
        global.fetch = vi.fn() as unknown as typeof fetch;
    });

    it('starts unauthenticated when /user/info returns 401', async () => {
        (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
            ok: false, status: 401, json: () => Promise.resolve({ authenticated: false }),
        });
        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.loading).toBe(false); });
        expect(result.current.user).toBeNull();
    });

    it('sets user after successful login', async () => {
        (global.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce({ ok: false, json: () => Promise.resolve({ authenticated: false }) })
            .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ success: true, expires_at: '2099-01-01T00:00:00Z' }) })
            .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ authenticated: true, username: 'alice', is_superuser: true }) });

        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.loading).toBe(false); });

        await act(async () => { await result.current.login('alice', 'pw'); });
        expect(result.current.user?.username).toBe('alice');
        expect(result.current.user?.isSuperuser).toBe(true);
    });

    it('clears user on logout', async () => {
        (global.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ authenticated: true, username: 'u' }) })
            .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ success: true }) });

        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.user?.username).toBe('u'); });
        await act(async () => { await result.current.logout(); });
        expect(result.current.user).toBeNull();
    });
});
