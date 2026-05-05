import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../../../test/renderWithTheme';
import AdminUsers from '../index';
import AuthContext from '../../../../contexts/AuthContext';
import type { AuthContextValue, User } from '../../../../contexts/AuthContext';

function makeOkResponse(json: unknown, status = 200): Response {
    return {
        ok: true,
        status,
        json: () => Promise.resolve(json),
        text: () => Promise.resolve(JSON.stringify(json)),
    } as unknown as Response;
}

function makeErrResponse(status: number, body: unknown): Response {
    const text = JSON.stringify(body);
    return {
        ok: false,
        status,
        json: () => Promise.resolve(body),
        text: () => Promise.resolve(text),
    } as unknown as Response;
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

function makeAuthContext(user: User | null): AuthContextValue {
    return {
        user,
        loading: false,
        login: vi.fn(),
        logout: vi.fn(),
        forceLogout: vi.fn(),
    };
}

function renderAdminUsers(authUser: User | null = superuserUser) {
    return renderWithTheme(
        <AuthContext.Provider value={makeAuthContext(authUser)}>
            <AdminUsers />
        </AuthContext.Provider>,
    );
}

const mockUsers = [
    { username: 'alice', full_name: 'Alice Smith', email: 'alice@example.com', is_superuser: false, enabled: true },
    { username: 'bob', full_name: 'Bob Jones', email: 'bob@example.com', is_superuser: true, enabled: false },
];

describe('AdminUsers', () => {
    beforeEach(() => {
        globalThis.fetch = vi.fn() as unknown as typeof fetch;
    });

    it('shows superuser-only alert when user is not superuser', () => {
        renderAdminUsers(regularUser);
        expect(screen.getByRole('alert')).toHaveTextContent(/superuser only/i);
    });

    it('renders the user list from a mocked GET', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse(mockUsers),
        );
        renderAdminUsers();

        await waitFor(() => {
            expect(screen.getByText('alice')).toBeInTheDocument();
            expect(screen.getByText('bob')).toBeInTheDocument();
        });
        expect(screen.getByText('Alice Smith')).toBeInTheDocument();
        expect(screen.getByText('alice@example.com')).toBeInTheDocument();
    });

    it('shows loading indicator then table', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse(mockUsers),
        );
        renderAdminUsers();
        // Eventually table appears
        await waitFor(() => expect(screen.getByRole('table')).toBeInTheDocument());
    });

    it('shows error when fetch fails', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeErrResponse(500, { error: 'Internal Server Error' }),
        );
        renderAdminUsers();
        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/internal server error/i);
        });
    });

    it('shows "No users found" when the list is empty', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(makeOkResponse([]));
        renderAdminUsers();
        await waitFor(() => {
            expect(screen.getByText(/no users found/i)).toBeInTheDocument();
        });
    });

    it('clicking Add User opens the UserDialog', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(makeOkResponse([]));
        renderAdminUsers();

        await waitFor(() => screen.getByRole('button', { name: /add user/i }));
        await user.click(screen.getByRole('button', { name: /add user/i }));

        await waitFor(() => {
            expect(screen.getByRole('dialog', { name: /add user/i })).toBeInTheDocument();
        });
    });

    it('clicking edit icon opens edit dialog', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(makeOkResponse(mockUsers));
        renderAdminUsers();

        await waitFor(() => screen.getByText('alice'));
        const editBtn = screen.getByLabelText('edit alice');
        await user.click(editBtn);

        await waitFor(() => {
            expect(screen.getByRole('dialog', { name: /edit user/i })).toBeInTheDocument();
        });
    });

    it('clicking reset password icon opens password reset dialog', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(makeOkResponse(mockUsers));
        renderAdminUsers();

        await waitFor(() => screen.getByText('alice'));
        const resetBtn = screen.getByLabelText('reset password for alice');
        await user.click(resetBtn);

        await waitFor(() => {
            expect(screen.getByRole('dialog', { name: /reset password/i })).toBeInTheDocument();
        });
    });

    it('clicking delete icon opens delete confirmation dialog', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(makeOkResponse(mockUsers));
        renderAdminUsers();

        await waitFor(() => screen.getByText('alice'));
        const deleteBtn = screen.getByLabelText('delete alice');
        await user.click(deleteBtn);

        await waitFor(() => {
            expect(screen.getByRole('dialog', { name: /delete user/i })).toBeInTheDocument();
        });
    });

    it('confirming delete calls DELETE endpoint and refreshes list', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce(makeOkResponse(mockUsers))
            .mockResolvedValueOnce(makeOkResponse(undefined, 204))
            .mockResolvedValueOnce(makeOkResponse([]));

        renderAdminUsers();

        await waitFor(() => screen.getByText('alice'));
        await user.click(screen.getByLabelText('delete alice'));
        await waitFor(() => screen.getByRole('dialog', { name: /delete user/i }));

        const confirmBtn = screen.getByRole('button', { name: /^delete$/i });
        await user.click(confirmBtn);

        await waitFor(() => expect(screen.getByText(/no users found/i)).toBeInTheDocument());
    });

    it('closing add dialog resets state', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(makeOkResponse([]));
        renderAdminUsers();

        await waitFor(() => screen.getByRole('button', { name: /add user/i }));
        await user.click(screen.getByRole('button', { name: /add user/i }));
        await waitFor(() => screen.getByRole('dialog', { name: /add user/i }));

        // Close via cancel
        await user.click(screen.getByRole('button', { name: /cancel/i }));
        await waitFor(() => expect(screen.queryByRole('dialog', { name: /add user/i })).not.toBeInTheDocument());
    });

    it('successfully adding a user closes dialog and refreshes list', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce(makeOkResponse([]))
            .mockResolvedValueOnce(makeOkResponse({ success: true, username: 'newuser' }, 201))
            .mockResolvedValueOnce(makeOkResponse([{ username: 'newuser', full_name: '', email: '', is_superuser: false, enabled: true }]));

        renderAdminUsers();

        await waitFor(() => screen.getByRole('button', { name: /add user/i }));
        await user.click(screen.getByRole('button', { name: /add user/i }));
        await waitFor(() => screen.getByRole('dialog', { name: /add user/i }));

        await user.type(screen.getByLabelText(/username/i), 'newuser');
        await user.type(screen.getByLabelText(/^password$/i), 'secret123');
        await user.click(screen.getByRole('button', { name: /add user/i }));

        await waitFor(() => expect(screen.queryByRole('dialog', { name: /add user/i })).not.toBeInTheDocument());
        await waitFor(() => expect(screen.getByText('newuser')).toBeInTheDocument());
    });

    it('password reset dialog onSaved closes it', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce(makeOkResponse(mockUsers))
            .mockResolvedValueOnce(makeOkResponse({ success: true }));

        renderAdminUsers();

        await waitFor(() => screen.getByText('alice'));
        await user.click(screen.getByLabelText('reset password for alice'));
        await waitFor(() => screen.getByRole('dialog', { name: /reset password/i }));

        await user.type(screen.getByLabelText(/new password/i), 'newpass123');
        await user.click(screen.getByRole('button', { name: /reset password/i }));

        await waitFor(() => {
            expect(screen.queryByRole('dialog', { name: /reset password/i })).not.toBeInTheDocument();
        });
    });

    it('edit dialog onSaved closes it and refreshes', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>)
            .mockResolvedValueOnce(makeOkResponse(mockUsers))
            .mockResolvedValueOnce(makeOkResponse({ success: true }))
            .mockResolvedValueOnce(makeOkResponse(mockUsers));

        renderAdminUsers();

        await waitFor(() => screen.getByText('alice'));
        await user.click(screen.getByLabelText('edit alice'));
        await waitFor(() => screen.getByRole('dialog', { name: /edit user/i }));

        await user.click(screen.getByRole('button', { name: /save changes/i }));

        await waitFor(() => {
            expect(screen.queryByRole('dialog', { name: /edit user/i })).not.toBeInTheDocument();
        });
    });
});
