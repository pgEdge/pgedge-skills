import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../../test/renderWithTheme';
import { axe } from '../../../test/axe';
import SettingsPanel from '../index';
import AuthContext from '../../../contexts/AuthContext';
import type { AuthContextValue, User } from '../../../contexts/AuthContext';

function makeAuthContext(user: User | null): AuthContextValue {
    return {
        user,
        loading: false,
        login: vi.fn(),
        logout: vi.fn(),
        forceLogout: vi.fn(),
    };
}

const superuserUser: User = {
    authenticated: true,
    username: 'admin',
    isSuperuser: true,
};

const regularUser: User = {
    authenticated: true,
    username: 'alice',
    isSuperuser: false,
};

function renderSettingsPanel(user: User | null = regularUser, onClose = vi.fn()) {
    return renderWithTheme(
        <AuthContext.Provider value={makeAuthContext(user)}>
            <SettingsPanel open={true} onClose={onClose} />
        </AuthContext.Provider>,
    );
}

function renderSettingsPanelWithAuth(user: User | null = superuserUser) {
    return renderSettingsPanel(user);
}

describe('SettingsPanel', () => {
    beforeEach(() => {
        globalThis.fetch = vi.fn() as unknown as typeof fetch;
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
            ok: true,
            status: 200,
            json: () => Promise.resolve([]),
            text: () => Promise.resolve('[]'),
        } as unknown as Response);
    });

    it('renders the title and the General nav item', () => {
        renderSettingsPanel();
        expect(screen.getByRole('dialog', { name: /settings/i })).toBeInTheDocument();
        expect(screen.getByRole('navigation', { name: 'settings sections' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /general/i })).toBeInTheDocument();
    });

    it('calls onClose when close button is clicked', async () => {
        const onClose = vi.fn();
        const user = userEvent.setup();
        renderSettingsPanel(regularUser, onClose);
        await user.click(screen.getByLabelText('close settings'));
        expect(onClose).toHaveBeenCalled();
    });

    it('selects the first item (Security/Users) on open', async () => {
        renderSettingsPanelWithAuth();
        await waitFor(() => {
            expect(screen.getByRole('button', { name: /users/i })).toBeInTheDocument();
        });
    });

    it('updates selection when General nav item is clicked', async () => {
        const user = userEvent.setup();
        renderSettingsPanel();
        const generalBtn = screen.getByRole('button', { name: /general/i });
        await user.click(generalBtn);
        expect(screen.getByText(/where settings go/i)).toBeInTheDocument();
    });

    it('does not render dialog content when closed', () => {
        renderWithTheme(
            <AuthContext.Provider value={makeAuthContext(regularUser)}>
                <SettingsPanel open={false} onClose={vi.fn()} />
            </AuthContext.Provider>,
        );
        expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });

    it('shows the Security/Users nav item', () => {
        renderSettingsPanel();
        expect(screen.getByRole('button', { name: /users/i })).toBeInTheDocument();
    });

    it('clicking Security/Users nav loads AdminUsers component for superuser', async () => {
        const user = userEvent.setup();
        renderSettingsPanelWithAuth();
        const usersBtn = screen.getByRole('button', { name: /users/i });
        await user.click(usersBtn);
        // AdminUsers renders a heading
        await waitFor(() => {
            expect(screen.getByRole('heading', { name: /users/i })).toBeInTheDocument();
        });
    });

    it('has no a11y violations', async () => {
        const { container } = renderSettingsPanel();
        expect(await axe(container)).toHaveNoViolations();
    });
});
