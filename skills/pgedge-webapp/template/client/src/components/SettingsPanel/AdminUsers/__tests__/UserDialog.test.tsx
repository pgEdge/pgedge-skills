import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../../../test/renderWithTheme';
import UserDialog from '../UserDialog';
import type { UserListItem } from '../index';

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

const existingUser: UserListItem = {
    username: 'alice',
    full_name: 'Alice Smith',
    email: 'alice@example.com',
    is_superuser: false,
    enabled: true,
};

describe('UserDialog — add mode', () => {
    beforeEach(() => {
        globalThis.fetch = vi.fn() as unknown as typeof fetch;
    });

    it('renders username, full name, email, and password fields', () => {
        renderWithTheme(
            <UserDialog mode="add" open={true} onClose={vi.fn()} onSaved={vi.fn()} />,
        );
        expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/full name/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/^password$/i)).toBeInTheDocument();
    });

    it('shows error when username is empty', async () => {
        const user = userEvent.setup();
        renderWithTheme(
            <UserDialog mode="add" open={true} onClose={vi.fn()} onSaved={vi.fn()} />,
        );
        await user.click(screen.getByRole('button', { name: /add user/i }));
        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/username is required/i);
        });
    });

    it('shows error when password is empty', async () => {
        const user = userEvent.setup();
        renderWithTheme(
            <UserDialog mode="add" open={true} onClose={vi.fn()} onSaved={vi.fn()} />,
        );
        await user.type(screen.getByLabelText(/username/i), 'newuser');
        await user.click(screen.getByRole('button', { name: /add user/i }));
        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/password is required/i);
        });
    });

    it('posts to /api/v1/users on valid add submit and calls onSaved', async () => {
        const onSaved = vi.fn();
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ success: true, username: 'newuser' }, 201),
        );

        renderWithTheme(
            <UserDialog mode="add" open={true} onClose={vi.fn()} onSaved={onSaved} />,
        );

        await user.type(screen.getByLabelText(/username/i), 'newuser');
        await user.type(screen.getByLabelText(/full name/i), 'New User');
        await user.type(screen.getByLabelText(/email/i), 'new@example.com');
        await user.type(screen.getByLabelText(/^password$/i), 'secret123');
        await user.click(screen.getByRole('button', { name: /add user/i }));

        await waitFor(() => expect(onSaved).toHaveBeenCalled());

        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0];
        expect(callArgs[0]).toBe('/api/v1/users');
        expect(callArgs[1].method).toBe('POST');
        const body = JSON.parse(callArgs[1].body as string);
        expect(body.username).toBe('newuser');
        expect(body.password).toBe('secret123');
    });

    it('surfaces error from server on add failure', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeErrResponse(409, { error: 'Username already exists' }),
        );

        renderWithTheme(
            <UserDialog mode="add" open={true} onClose={vi.fn()} onSaved={vi.fn()} />,
        );

        await user.type(screen.getByLabelText(/username/i), 'alice');
        await user.type(screen.getByLabelText(/^password$/i), 'secret123');
        await user.click(screen.getByRole('button', { name: /add user/i }));

        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/username already exists/i);
        });
    });

    it('calls onClose when cancel is clicked', async () => {
        const onClose = vi.fn();
        const user = userEvent.setup();
        renderWithTheme(
            <UserDialog mode="add" open={true} onClose={onClose} onSaved={vi.fn()} />,
        );
        await user.click(screen.getByRole('button', { name: /cancel/i }));
        expect(onClose).toHaveBeenCalled();
    });
});

describe('UserDialog — edit mode', () => {
    beforeEach(() => {
        globalThis.fetch = vi.fn() as unknown as typeof fetch;
    });

    it('does NOT show password field in edit mode', () => {
        renderWithTheme(
            <UserDialog
                mode="edit"
                initialUser={existingUser}
                open={true}
                onClose={vi.fn()}
                onSaved={vi.fn()}
            />,
        );
        expect(screen.queryByLabelText(/^password$/i)).not.toBeInTheDocument();
    });

    it('shows Enabled switch in edit mode', () => {
        renderWithTheme(
            <UserDialog
                mode="edit"
                initialUser={existingUser}
                open={true}
                onClose={vi.fn()}
                onSaved={vi.fn()}
            />,
        );
        expect(screen.getByLabelText(/enabled/i)).toBeInTheDocument();
    });

    it('pre-fills fields from initialUser', () => {
        renderWithTheme(
            <UserDialog
                mode="edit"
                initialUser={existingUser}
                open={true}
                onClose={vi.fn()}
                onSaved={vi.fn()}
            />,
        );
        expect(screen.getByDisplayValue('alice')).toBeInTheDocument();
        expect(screen.getByDisplayValue('Alice Smith')).toBeInTheDocument();
        expect(screen.getByDisplayValue('alice@example.com')).toBeInTheDocument();
    });

    it('sends PATCH to /api/v1/users/{username} on save', async () => {
        const onSaved = vi.fn();
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ success: true }),
        );

        renderWithTheme(
            <UserDialog
                mode="edit"
                initialUser={existingUser}
                open={true}
                onClose={vi.fn()}
                onSaved={onSaved}
            />,
        );

        // Clear and retype full name
        const fullNameInput = screen.getByLabelText(/full name/i);
        await user.clear(fullNameInput);
        await user.type(fullNameInput, 'Alice Updated');
        await user.click(screen.getByRole('button', { name: /save changes/i }));

        await waitFor(() => expect(onSaved).toHaveBeenCalled());

        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0];
        expect(callArgs[0]).toBe('/api/v1/users/alice');
        expect(callArgs[1].method).toBe('PATCH');
        const body = JSON.parse(callArgs[1].body as string);
        expect(body.full_name).toBe('Alice Updated');
    });

    it('surfaces error from server on edit failure', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeErrResponse(500, { error: 'Server error' }),
        );

        renderWithTheme(
            <UserDialog
                mode="edit"
                initialUser={existingUser}
                open={true}
                onClose={vi.fn()}
                onSaved={vi.fn()}
            />,
        );

        await user.click(screen.getByRole('button', { name: /save changes/i }));

        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/server error/i);
        });
    });

    it('toggles Is Superuser switch in edit mode', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ success: true }),
        );
        const onSaved = vi.fn();

        renderWithTheme(
            <UserDialog
                mode="edit"
                initialUser={existingUser}
                open={true}
                onClose={vi.fn()}
                onSaved={onSaved}
            />,
        );

        // Toggle Is Superuser switch
        await user.click(screen.getByLabelText(/is superuser/i));
        await user.click(screen.getByRole('button', { name: /save changes/i }));

        await waitFor(() => expect(onSaved).toHaveBeenCalled());

        const body = JSON.parse(
            (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1].body as string,
        );
        expect(body.is_superuser).toBe(true);
    });

    it('toggles Enabled switch in edit mode', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ success: true }),
        );
        const onSaved = vi.fn();

        renderWithTheme(
            <UserDialog
                mode="edit"
                initialUser={existingUser}
                open={true}
                onClose={vi.fn()}
                onSaved={onSaved}
            />,
        );

        // Toggle Enabled switch (currently true — switch it off)
        await user.click(screen.getByLabelText(/^enabled$/i));
        await user.click(screen.getByRole('button', { name: /save changes/i }));

        await waitFor(() => expect(onSaved).toHaveBeenCalled());

        const body = JSON.parse(
            (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1].body as string,
        );
        expect(body.enabled).toBe(false);
    });
});
