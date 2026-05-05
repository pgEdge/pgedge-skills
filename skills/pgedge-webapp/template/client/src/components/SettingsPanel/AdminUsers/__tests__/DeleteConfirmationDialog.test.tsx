import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../../../test/renderWithTheme';
import DeleteConfirmationDialog from '../DeleteConfirmationDialog';

describe('DeleteConfirmationDialog', () => {
    beforeEach(() => {
        globalThis.fetch = vi.fn() as unknown as typeof fetch;
    });

    it('renders with username in content', () => {
        renderWithTheme(
            <DeleteConfirmationDialog
                username="alice"
                open={true}
                onClose={vi.fn()}
                onConfirm={vi.fn()}
            />,
        );
        expect(screen.getByRole('dialog', { name: /delete user/i })).toBeInTheDocument();
        expect(screen.getByText(/alice/i)).toBeInTheDocument();
    });

    it('does not render when closed', () => {
        renderWithTheme(
            <DeleteConfirmationDialog
                username="alice"
                open={false}
                onClose={vi.fn()}
                onConfirm={vi.fn()}
            />,
        );
        expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });

    it('calls onConfirm with username when delete is clicked', async () => {
        const onConfirm = vi.fn().mockResolvedValue(undefined);
        const user = userEvent.setup();
        renderWithTheme(
            <DeleteConfirmationDialog
                username="alice"
                open={true}
                onClose={vi.fn()}
                onConfirm={onConfirm}
            />,
        );
        await user.click(screen.getByRole('button', { name: /^delete$/i }));
        await waitFor(() => {
            expect(onConfirm).toHaveBeenCalledWith('alice');
        });
    });

    it('calls onClose when cancel is clicked', async () => {
        const onClose = vi.fn();
        const user = userEvent.setup();
        renderWithTheme(
            <DeleteConfirmationDialog
                username="alice"
                open={true}
                onClose={onClose}
                onConfirm={vi.fn()}
            />,
        );
        await user.click(screen.getByRole('button', { name: /cancel/i }));
        expect(onClose).toHaveBeenCalled();
    });

    it('shows error when onConfirm throws', async () => {
        const onConfirm = vi.fn().mockRejectedValue(new Error('Permission denied'));
        const user = userEvent.setup();
        renderWithTheme(
            <DeleteConfirmationDialog
                username="alice"
                open={true}
                onClose={vi.fn()}
                onConfirm={onConfirm}
            />,
        );
        await user.click(screen.getByRole('button', { name: /^delete$/i }));
        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/permission denied/i);
        });
    });

    it('shows generic error when onConfirm throws non-Error', async () => {
        const onConfirm = vi.fn().mockRejectedValue('unknown error');
        const user = userEvent.setup();
        renderWithTheme(
            <DeleteConfirmationDialog
                username="alice"
                open={true}
                onClose={vi.fn()}
                onConfirm={onConfirm}
            />,
        );
        await user.click(screen.getByRole('button', { name: /^delete$/i }));
        await waitFor(() => {
            expect(screen.getByRole('alert')).toHaveTextContent(/failed to delete user/i);
        });
    });

    it('clears error and submitting state when cancel is clicked after error', async () => {
        const onConfirm = vi.fn().mockRejectedValue(new Error('Permission denied'));
        const onClose = vi.fn();
        const user = userEvent.setup();
        renderWithTheme(
            <DeleteConfirmationDialog
                username="alice"
                open={true}
                onClose={onClose}
                onConfirm={onConfirm}
            />,
        );
        // Trigger error
        await user.click(screen.getByRole('button', { name: /^delete$/i }));
        await waitFor(() => screen.getByRole('alert'));

        // Now cancel
        await user.click(screen.getByRole('button', { name: /cancel/i }));
        expect(onClose).toHaveBeenCalled();
    });
});
