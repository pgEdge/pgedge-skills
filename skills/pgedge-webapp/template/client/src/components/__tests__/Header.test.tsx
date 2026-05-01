import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../test/renderWithTheme';
import { axe } from '../../test/axe';
import Header from '../Header';
import { AuthProvider } from '../../contexts/AuthContext';
import { createTheme, ThemeProvider } from '@mui/material/styles';

function makeOkResponse(json: unknown): Response {
    return {
        ok: true,
        status: 200,
        json: () => Promise.resolve(json),
        text: () => Promise.resolve(JSON.stringify(json)),
    } as unknown as Response;
}

const renderHeader = (props: Partial<{ onToggleTheme: () => void }> = {}) => {
    (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
        makeOkResponse({ authenticated: true, username: 'tester' }),
    );
    return renderWithTheme(
        <AuthProvider><Header onToggleTheme={vi.fn()} {...props} /></AuthProvider>
    );
};

const renderHeaderWithUsername = (username: string, props: Partial<{ onToggleTheme: () => void }> = {}) => {
    (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
        makeOkResponse({ authenticated: true, username }),
    );
    return renderWithTheme(
        <AuthProvider><Header onToggleTheme={vi.fn()} {...props} /></AuthProvider>
    );
};

describe('Header', () => {
    beforeEach(() => { globalThis.fetch = vi.fn() as unknown as typeof fetch; });

    it('shows the project name', async () => {
        renderHeader();
        await waitFor(() => expect(screen.getByText('<PROJECT_NAME>')).toBeInTheDocument());
    });

    it('exposes theme/settings/help/user-menu buttons with aria-labels', async () => {
        renderHeader();
        await waitFor(() => screen.getByLabelText('toggle theme'));
        expect(screen.getByLabelText('open settings')).toBeInTheDocument();
        expect(screen.getByLabelText('open help')).toBeInTheDocument();
        expect(screen.getByLabelText('user menu')).toBeInTheDocument();
    });

    it('invokes onToggleTheme when clicked', async () => {
        const onToggleTheme = vi.fn();
        const user = userEvent.setup();
        renderHeader({ onToggleTheme });
        await user.click(await screen.findByLabelText('toggle theme'));
        expect(onToggleTheme).toHaveBeenCalledTimes(1);
    });

    it('opens the user menu and signs out', async () => {
        const user = userEvent.setup();
        renderHeader();
        await user.click(await screen.findByLabelText('user menu'));
        const signOut = await screen.findByText(/Sign out/i);
        await user.click(signOut);
    });

    it('opens the help panel when help button is clicked', async () => {
        const user = userEvent.setup();
        renderHeader();
        await waitFor(() => screen.getByLabelText('open help'));
        await user.click(screen.getByLabelText('open help'));
        // Help panel drawer should appear
        await waitFor(() => expect(screen.getByLabelText('close help')).toBeInTheDocument());
    });

    it('closes the help panel when close is clicked', async () => {
        const user = userEvent.setup();
        renderHeader();
        await waitFor(() => screen.getByLabelText('open help'));
        await user.click(screen.getByLabelText('open help'));
        await waitFor(() => screen.getByLabelText('close help'));
        await user.click(screen.getByLabelText('close help'));
        await waitFor(() => expect(screen.queryByLabelText('close help')).not.toBeInTheDocument());
    });

    it('opens the settings panel when settings button is clicked', async () => {
        const user = userEvent.setup();
        renderHeader();
        await waitFor(() => screen.getByLabelText('open settings'));
        await user.click(screen.getByLabelText('open settings'));
        // Settings dialog should appear
        await waitFor(() => expect(screen.getByRole('dialog', { name: /settings/i })).toBeInTheDocument());
    });

    it('closes the settings panel when close is clicked', async () => {
        const user = userEvent.setup();
        renderHeader();
        await waitFor(() => screen.getByLabelText('open settings'));
        await user.click(screen.getByLabelText('open settings'));
        await waitFor(() => screen.getByLabelText('close settings'));
        await user.click(screen.getByLabelText('close settings'));
        await waitFor(() => expect(screen.queryByRole('dialog', { name: /settings/i })).not.toBeInTheDocument());
    });

    it('shows single-char initials for single-word username', async () => {
        renderHeaderWithUsername('alice');
        await waitFor(() => screen.getByLabelText('user menu'));
        // Avatar content is first char uppercased
        expect(screen.getByText('A')).toBeInTheDocument();
    });

    it('shows two-char initials for multi-word username', async () => {
        renderHeaderWithUsername('Alice Smith');
        await waitFor(() => screen.getByLabelText('user menu'));
        expect(screen.getByText('AS')).toBeInTheDocument();
    });

    it('shows ? for undefined username', async () => {
        // Mock checkAuth returning authenticated=true but no username
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ authenticated: true, username: '' }),
        );
        renderWithTheme(
            <AuthProvider><Header onToggleTheme={vi.fn()} /></AuthProvider>
        );
        await waitFor(() => screen.getByLabelText('user menu'));
        expect(screen.getByText('?')).toBeInTheDocument();
    });

    it('has no a11y violations', async () => {
        const { container } = renderHeader();
        await waitFor(() => screen.getByLabelText('toggle theme'));
        expect(await axe(container)).toHaveNoViolations();
    });
});
