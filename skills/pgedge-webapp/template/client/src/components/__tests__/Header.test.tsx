import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../test/renderWithTheme';
import { axe } from '../../test/axe';
import Header from '../Header';
import { AuthProvider } from '../../contexts/AuthContext';

const renderHeader = (props: Partial<{ onToggleTheme: () => void }> = {}) => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true, json: () => Promise.resolve({ authenticated: true, username: 'tester' }),
    });
    return renderWithTheme(
        <AuthProvider><Header onToggleTheme={vi.fn()} {...props} /></AuthProvider>
    );
};

describe('Header', () => {
    beforeEach(() => { global.fetch = vi.fn() as unknown as typeof fetch; });

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

    it('has no a11y violations', async () => {
        const { container } = renderHeader();
        await waitFor(() => screen.getByLabelText('toggle theme'));
        expect(await axe(container)).toHaveNoViolations();
    });
});
