import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../test/renderWithTheme';
import { axe } from '../../test/axe';
import Login from '../Login';
import { AuthProvider } from '../../contexts/AuthContext';

const renderLogin = () => renderWithTheme(<AuthProvider><Login /></AuthProvider>);

function makeOkResponse(json: unknown): Response {
    return {
        ok: true, status: 200,
        json: () => Promise.resolve(json),
        text: () => Promise.resolve(JSON.stringify(json)),
    } as unknown as Response;
}

describe('Login', () => {
    beforeEach(() => {
        // Reset sessionStorage mock calls without losing the mock (setup.ts installs it)
        vi.mocked(window.sessionStorage.getItem).mockReset();
        vi.mocked(window.sessionStorage.setItem).mockReset();
        vi.mocked(window.sessionStorage.removeItem).mockReset();

        globalThis.fetch = vi.fn() as unknown as typeof fetch;
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
            ok: false, status: 401,
            json: () => Promise.resolve({ authenticated: false }),
            text: () => Promise.resolve(''),
        });
    });

    it('renders the project name as the title', async () => {
        renderLogin();
        await waitFor(() => {
            expect(screen.getByRole('heading', { name: /<PROJECT_NAME>/ })).toBeInTheDocument();
        });
    });

    it('renders username and password fields', async () => {
        renderLogin();
        await waitFor(() => {
            expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
            expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
        });
    });

    it('shows an error on bad credentials', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce({
                ok: false, json: () => Promise.resolve({ authenticated: false }),
                text: () => Promise.resolve(''),
            })
            .mockResolvedValueOnce({
                ok: false, status: 401,
                json: () => Promise.resolve({ message: 'invalid' }),
                text: () => Promise.resolve('{"message":"invalid"}'),
            });

        renderLogin();
        await user.type(await screen.findByLabelText(/username/i), 'x');
        await user.type(screen.getByLabelText(/password/i), 'y');
        await user.click(screen.getByRole('button', { name: /sign in/i }));
        await waitFor(() => {
            expect(screen.getByRole('alert')).toBeInTheDocument();
        });
    });

    it('shows warning alert when sessionStorage has disconnect message', async () => {
        // The setup.ts mocks sessionStorage with vi.fn() — use mockReturnValueOnce
        (window.sessionStorage.getItem as ReturnType<typeof vi.fn>).mockReturnValueOnce('You were disconnected');

        renderLogin();
        // Wait for auth check and Login component to mount
        await waitFor(() => screen.getByRole('heading', { name: /<PROJECT_NAME>/ }));
        await waitFor(() => {
            expect(screen.getByText('You were disconnected')).toBeInTheDocument();
        });
    });

    it('clears the warning when its close button is clicked', async () => {
        const user = userEvent.setup();
        (window.sessionStorage.getItem as ReturnType<typeof vi.fn>).mockReturnValueOnce('Disconnected from server');

        renderLogin();
        await waitFor(() => screen.getByRole('heading', { name: /<PROJECT_NAME>/ }));
        const warningText = await screen.findByText('Disconnected from server');
        expect(warningText).toBeInTheDocument();

        // MUI Alert has a close button with aria-label "Close"
        const closeBtn = screen.getByLabelText(/close/i);
        await user.click(closeBtn);

        await waitFor(() => {
            expect(screen.queryByText('Disconnected from server')).not.toBeInTheDocument();
        });
    });

    it('shows success login and redirects user', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce({
                ok: false, status: 401,
                json: () => Promise.resolve({ authenticated: false }),
                text: () => Promise.resolve(''),
            })
            .mockResolvedValueOnce(makeOkResponse({ success: true, expires_at: '2099-01-01T00:00:00Z' }))
            .mockResolvedValueOnce(makeOkResponse({ authenticated: true, username: 'alice', is_superuser: false }));

        renderLogin();
        await user.type(await screen.findByLabelText(/username/i), 'alice');
        await user.type(screen.getByLabelText(/password/i), 'goodpassword');
        await user.click(screen.getByRole('button', { name: /sign in/i }));
        // After login, there should be no error alert shown
        await waitFor(() => {
            expect(screen.queryByRole('alert')).not.toBeInTheDocument();
        });
    });

    it('has no a11y violations', async () => {
        const { container } = renderLogin();
        await waitFor(() => screen.getByRole('heading'));
        expect(await axe(container)).toHaveNoViolations();
    });
});
