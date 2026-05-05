import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../test/renderWithTheme';
import ChangePasswordDialog from '../Header/ChangePasswordDialog';

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

describe('ChangePasswordDialog', () => {
    beforeEach(() => {
        globalThis.fetch = vi.fn() as unknown as typeof fetch;
    });

    it('renders three password fields', () => {
        renderWithTheme(<ChangePasswordDialog open={true} onClose={vi.fn()} />);
        expect(screen.getByLabelText(/current password/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/new password/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/confirm password/i)).toBeInTheDocument();
    });

    it('does not render when closed', () => {
        renderWithTheme(<ChangePasswordDialog open={false} onClose={vi.fn()} />);
        expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });

    it('shows error when passwords do not match', async () => {
        const user = userEvent.setup();
        renderWithTheme(<ChangePasswordDialog open={true} onClose={vi.fn()} />);

        await user.type(screen.getByLabelText(/current password/i), 'old123');
        await user.type(screen.getByLabelText(/new password/i), 'newpass1');
        await user.type(screen.getByLabelText(/confirm password/i), 'newpass2');
        await user.click(screen.getByRole('button', { name: /change password/i }));

        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/do not match/i);
        });
    });

    it('shows error when current password is empty', async () => {
        const user = userEvent.setup();
        renderWithTheme(<ChangePasswordDialog open={true} onClose={vi.fn()} />);

        await user.type(screen.getByLabelText(/new password/i), 'newpass1');
        await user.type(screen.getByLabelText(/confirm password/i), 'newpass1');
        await user.click(screen.getByRole('button', { name: /change password/i }));

        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/current password is required/i);
        });
    });

    it('shows error when new password is empty', async () => {
        const user = userEvent.setup();
        renderWithTheme(<ChangePasswordDialog open={true} onClose={vi.fn()} />);

        await user.type(screen.getByLabelText(/current password/i), 'old123');
        await user.click(screen.getByRole('button', { name: /change password/i }));

        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/new password is required/i);
        });
    });

    it('posts to /api/v1/user/password on valid submit and calls onClose', async () => {
        const onClose = vi.fn();
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeOkResponse({ success: true }),
        );

        renderWithTheme(<ChangePasswordDialog open={true} onClose={onClose} />);

        await user.type(screen.getByLabelText(/current password/i), 'old123');
        await user.type(screen.getByLabelText(/new password/i), 'newpass1');
        await user.type(screen.getByLabelText(/confirm password/i), 'newpass1');
        await user.click(screen.getByRole('button', { name: /change password/i }));

        await waitFor(() => expect(onClose).toHaveBeenCalled());

        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0];
        expect(callArgs[0]).toBe('/api/v1/user/password');
        expect(callArgs[1].method).toBe('POST');
        const body = JSON.parse(callArgs[1].body as string);
        expect(body.current_password).toBe('old123');
        expect(body.new_password).toBe('newpass1');
    });

    it('surfaces error message from server on failure', async () => {
        const user = userEvent.setup();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeErrResponse(400, { error: 'Current password is incorrect' }),
        );

        renderWithTheme(<ChangePasswordDialog open={true} onClose={vi.fn()} />);

        await user.type(screen.getByLabelText(/current password/i), 'wrongpass');
        await user.type(screen.getByLabelText(/new password/i), 'newpass1');
        await user.type(screen.getByLabelText(/confirm password/i), 'newpass1');
        await user.click(screen.getByRole('button', { name: /change password/i }));

        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/current password is incorrect/i);
        });
    });

    it('calls onClose when cancel is clicked', async () => {
        const onClose = vi.fn();
        const user = userEvent.setup();
        renderWithTheme(<ChangePasswordDialog open={true} onClose={onClose} />);

        await user.click(screen.getByRole('button', { name: /cancel/i }));
        expect(onClose).toHaveBeenCalled();
    });
});
