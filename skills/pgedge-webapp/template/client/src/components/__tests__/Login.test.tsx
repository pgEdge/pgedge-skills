import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../test/renderWithTheme';
import { axe } from '../../test/axe';
import Login from '../Login';
import { AuthProvider } from '../../contexts/AuthContext';

const renderLogin = () => renderWithTheme(<AuthProvider><Login /></AuthProvider>);

describe('Login', () => {
    beforeEach(() => {
        global.fetch = vi.fn() as unknown as typeof fetch;
        (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
            ok: false, status: 401, json: () => Promise.resolve({ authenticated: false }),
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
        (global.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce({ ok: false, json: () => Promise.resolve({ authenticated: false }) })
            .mockResolvedValueOnce({ ok: false, status: 401, json: () => Promise.resolve({ message: 'invalid' }) });

        renderLogin();
        await user.type(await screen.findByLabelText(/username/i), 'x');
        await user.type(screen.getByLabelText(/password/i), 'y');
        await user.click(screen.getByRole('button', { name: /sign in/i }));
        await waitFor(() => {
            expect(screen.getByRole('alert')).toBeInTheDocument();
        });
    });

    it('has no a11y violations', async () => {
        const { container } = renderLogin();
        await waitFor(() => screen.getByRole('heading'));
        expect(await axe(container)).toHaveNoViolations();
    });
});
