import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from './test/renderWithTheme';
import App from './App';

function makeOkResponse(json: unknown): Response {
    return {
        ok: true, status: 200,
        json: () => Promise.resolve(json),
        text: () => Promise.resolve(JSON.stringify(json)),
    } as unknown as Response;
}

describe('App', () => {
    beforeEach(() => {
        globalThis.fetch = vi.fn() as unknown as typeof fetch;
        vi.mocked(window.localStorage.getItem).mockReturnValue(null);
    });

    it('shows loading spinner initially', async () => {
        // Never resolve so we can catch the loading state
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockReturnValue(new Promise(() => undefined));
        renderWithTheme(<App />);
        expect(screen.getByLabelText('Loading application')).toBeInTheDocument();
    });

    it('shows login page when not authenticated', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
            ok: false, status: 401,
            json: () => Promise.resolve({ authenticated: false }),
            text: () => Promise.resolve(''),
        });
        renderWithTheme(<App />);
        await waitFor(() => expect(screen.getByRole('heading', { name: /<PROJECT_NAME>/ })).toBeInTheDocument());
    });

    it('shows main content when authenticated', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ authenticated: true, username: 'alice' }),
        );
        renderWithTheme(<App />);
        await waitFor(() => expect(screen.getByLabelText('toggle theme')).toBeInTheDocument());
    });

    it('toggles theme from light to dark via onToggleTheme', async () => {
        const user = userEvent.setup();
        vi.mocked(window.localStorage.getItem).mockReturnValue('light');
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ authenticated: true, username: 'alice' }),
        );
        renderWithTheme(<App />);
        await waitFor(() => screen.getByLabelText('toggle theme'));
        // Click once: light → dark (should now show LightModeIcon)
        await user.click(screen.getByLabelText('toggle theme'));
        // localStorage.setItem should be called with 'dark'
        expect(vi.mocked(window.localStorage.setItem)).toHaveBeenCalledWith('theme-mode', 'dark');
    });

    it('starts in dark mode if localStorage has dark', async () => {
        vi.mocked(window.localStorage.getItem).mockReturnValue('dark');
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ authenticated: true, username: 'alice' }),
        );
        renderWithTheme(<App />);
        await waitFor(() => screen.getByLabelText('toggle theme'));
        // In dark mode, toggle should go back to light
        const toggleBtn = screen.getByLabelText('toggle theme');
        const user = userEvent.setup();
        await user.click(toggleBtn);
        expect(vi.mocked(window.localStorage.setItem)).toHaveBeenCalledWith('theme-mode', 'light');
    });
});
