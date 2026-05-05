import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../../../test/renderWithTheme';
import AdminPasswordResetDialog from '../AdminPasswordResetDialog';

function makeOkResponse(json: unknown): Response {
    return {
        ok: true,
        status: 200,
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

describe('AdminPasswordResetDialog', () => {
    beforeEach(() => {
        globalThis.fetch = vi.fn() as unknown as typeof fetch;
    });

    it('renders with username in title', () => {
        renderWithTheme(
            <AdminPasswordResetDialog
                username="alice"
                open={true}
                onClose={vi.fn()}
                onSaved={vi.fn()}
            />,
        );
        expect(screen.getByRole('dialog', { name: /reset password: alice/i })).toBeInTheDocument();
    });

    it('does not render when closed', () => {
        renderWithTheme(
            <AdminPasswordResetDialog
                username="alice"
                open={false}
                onClose={vi.fn()}
                onSaved={vi.fn()}
            />,
        );
        expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });

    it('shows error when password is empty', async () => {
        const user = userEvent.setup();
        renderWithTheme(
            <AdminPasswordResetDialog
                username="alice"
                open={true}
                onClose={vi.fn()}
                onSaved={vi.fn()}
            />,
        );
        await user.click(screen.getByRole('button', { name: /reset password/i }));
        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/new password is required/i);
        });
    });

    it('posts to /api/v1/users/{username}/password on valid submit', async () => {
        const onSaved = vi.fn();
        const onClose = vi.fn();
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ success: true }),
        );

        renderWithTheme(
            <AdminPasswordResetDialog
                username="alice"
                open={true}
                onClose={onClose}
                onSaved={onSaved}
            />,
        );

        await user.type(screen.getByLabelText(/new password/i), 'newSecret123');
        await user.click(screen.getByRole('button', { name: /reset password/i }));

        await waitFor(() => {
            expect(onSaved).toHaveBeenCalled();
            expect(onClose).toHaveBeenCalled();
        });

        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0];
        expect(callArgs[0]).toBe('/api/v1/users/alice/password');
        expect(callArgs[1].method).toBe('POST');
        const body = JSON.parse(callArgs[1].body as string);
        expect(body.new_password).toBe('newSecret123');
    });

    it('surfaces error from server on failure', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeErrResponse(500, { error: 'Server failure' }),
        );

        renderWithTheme(
            <AdminPasswordResetDialog
                username="alice"
                open={true}
                onClose={vi.fn()}
                onSaved={vi.fn()}
            />,
        );

        await user.type(screen.getByLabelText(/new password/i), 'newSecret123');
        await user.click(screen.getByRole('button', { name: /reset password/i }));

        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/server failure/i);
        });
    });

    it('calls onClose when cancel is clicked', async () => {
        const onClose = vi.fn();
        const user = userEvent.setup();
        renderWithTheme(
            <AdminPasswordResetDialog
                username="alice"
                open={true}
                onClose={onClose}
                onSaved={vi.fn()}
            />,
        );
        await user.click(screen.getByRole('button', { name: /cancel/i }));
        expect(onClose).toHaveBeenCalled();
    });

    it('resets state when reopened', async () => {
        const { rerender } = renderWithTheme(
            <AdminPasswordResetDialog
                username="alice"
                open={false}
                onClose={vi.fn()}
                onSaved={vi.fn()}
            />,
        );
        rerender(
            <AdminPasswordResetDialog
                username="alice"
                open={true}
                onClose={vi.fn()}
                onSaved={vi.fn()}
            />,
        );
        // Should have an empty password field
        expect(screen.getByLabelText(/new password/i)).toHaveValue('');
    });
});
