import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { AuthProvider } from '../AuthContext';
import { useAuth } from '../useAuth';

const wrapper = ({ children }: { children: React.ReactNode }) => (
    <AuthProvider>{children}</AuthProvider>
);

function makeOkResponse(json: unknown): Response {
    return {
        ok: true,
        status: 200,
        json: () => Promise.resolve(json),
        text: () => Promise.resolve(JSON.stringify(json)),
    } as unknown as Response;
}


describe('AuthContext', () => {
    beforeEach(() => {
        globalThis.fetch = vi.fn() as unknown as typeof fetch;
    });

    it('starts unauthenticated when /user/info returns 401', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
            ok: false, status: 401, json: () => Promise.resolve({ authenticated: false }),
            text: () => Promise.resolve(''),
        });
        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.loading).toBe(false); });
        expect(result.current.user).toBeNull();
    });

    it('sets user when /user/info returns authenticated', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ authenticated: true, username: 'bob', is_superuser: false }),
        );
        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.loading).toBe(false); });
        expect(result.current.user?.username).toBe('bob');
        expect(result.current.user?.isSuperuser).toBe(false);
    });

    it('sets user to null when /user/info returns authenticated=false', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ authenticated: false, username: '' }),
        );
        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.loading).toBe(false); });
        expect(result.current.user).toBeNull();
    });

    it('sets user after successful login', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce({ ok: false, status: 401, json: () => Promise.resolve({ authenticated: false }), text: () => Promise.resolve('') })
            .mockResolvedValueOnce(makeOkResponse({ success: true, expires_at: '2099-01-01T00:00:00Z' }))
            .mockResolvedValueOnce(makeOkResponse({ authenticated: true, username: 'alice', is_superuser: true }));

        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.loading).toBe(false); });

        await act(async () => { await result.current.login('alice', 'pw'); });
        expect(result.current.user?.username).toBe('alice');
        expect(result.current.user?.isSuperuser).toBe(true);
    });

    it('throws when login returns success=false', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce({ ok: false, status: 401, json: () => Promise.resolve({ authenticated: false }), text: () => Promise.resolve('') })
            .mockResolvedValueOnce(makeOkResponse({ success: false, message: 'Invalid credentials' }));

        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.loading).toBe(false); });

        await expect(
            act(async () => { await result.current.login('alice', 'wrong'); }),
        ).rejects.toThrow('Invalid credentials');
    });

    it('throws with default message when login success=false and no message', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce({ ok: false, status: 401, json: () => Promise.resolve({ authenticated: false }), text: () => Promise.resolve('') })
            .mockResolvedValueOnce(makeOkResponse({ success: false }));

        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.loading).toBe(false); });

        await expect(
            act(async () => { await result.current.login('alice', 'wrong'); }),
        ).rejects.toThrow('Authentication failed');
    });

    it('sets user with only username when post-login /user/info fetch fails', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce({ ok: false, status: 401, json: () => Promise.resolve({ authenticated: false }), text: () => Promise.resolve('') })
            .mockResolvedValueOnce(makeOkResponse({ success: true, expires_at: '2099-01-01T00:00:00Z' }))
            .mockRejectedValueOnce(new TypeError('network error'));

        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.loading).toBe(false); });

        await act(async () => { await result.current.login('carol', 'pw'); });
        expect(result.current.user?.username).toBe('carol');
        expect(result.current.user?.authenticated).toBe(true);
    });

    it('clears user on logout', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce(makeOkResponse({ authenticated: true, username: 'u' }))
            .mockResolvedValueOnce(makeOkResponse({ success: true }));

        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.user?.username).toBe('u'); });
        await act(async () => { await result.current.logout(); });
        expect(result.current.user).toBeNull();
    });

    it('still clears user when logout POST fails', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce(makeOkResponse({ authenticated: true, username: 'u' }))
            .mockRejectedValueOnce(new TypeError('network error'));

        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.user?.username).toBe('u'); });
        await act(async () => { await result.current.logout(); });
        expect(result.current.user).toBeNull();
    });

    it('forceLogout clears user', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce(makeOkResponse({ authenticated: true, username: 'u' }))
            .mockResolvedValueOnce(makeOkResponse({ success: true }));

        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.user?.username).toBe('u'); });
        await act(async () => { await result.current.forceLogout(); });
        expect(result.current.user).toBeNull();
    });

    it('forceLogout clears user even when logout POST fails', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce(makeOkResponse({ authenticated: true, username: 'u' }))
            .mockRejectedValueOnce(new TypeError('network error'));

        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.user?.username).toBe('u'); });
        await act(async () => { await result.current.forceLogout(); });
        expect(result.current.user).toBeNull();
    });

    it('checkAuth error path: swallows errors and ends with user=null, loading=false', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockRejectedValue(
            new TypeError('network error'),
        );
        const { result } = renderHook(() => useAuth(), { wrapper });
        await waitFor(() => { expect(result.current.loading).toBe(false); });
        expect(result.current.user).toBeNull();
    });
});

describe('useAuth outside provider', () => {
    it('throws when used outside AuthProvider', () => {
        // Suppress error boundary output
        const spy = vi.spyOn(console, 'error').mockImplementation(() => undefined);
        expect(() => {
            renderHook(() => useAuth());
        }).toThrow('useAuth must be used within AuthProvider');
        spy.mockRestore();
    });
});
