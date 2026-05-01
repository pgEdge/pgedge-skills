import { describe, it, expect, vi } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../../test/renderWithTheme';
import { axe } from '../../../test/axe';
import EllieChatPanel from '../index';

describe('EllieChatPanel', () => {
    it('renders the chat shell', () => {
        renderWithTheme(<EllieChatPanel open={true} onClose={vi.fn()} />);
        expect(screen.getByRole('heading', { name: /Ellie/i })).toBeInTheDocument();
        expect(screen.getByLabelText(/Message Ellie/)).toBeInTheDocument();
    });

    it('appends a user message on send', async () => {
        const user = userEvent.setup();
        renderWithTheme(<EllieChatPanel open={true} onClose={vi.fn()} />);
        const input = screen.getByLabelText(/Message Ellie/);
        await user.type(input, 'hello');
        await user.click(screen.getByLabelText('send message'));
        expect(screen.getByText('hello')).toBeInTheDocument();
    });

    it('does not send when input is empty or whitespace', async () => {
        const user = userEvent.setup();
        renderWithTheme(<EllieChatPanel open={true} onClose={vi.fn()} />);
        // Send button should be disabled when input is empty
        const sendBtn = screen.getByLabelText('send message');
        expect(sendBtn).toBeDisabled();

        // Type only whitespace - button should still be disabled
        const input = screen.getByLabelText(/Message Ellie/);
        await user.type(input, '   ');
        expect(sendBtn).toBeDisabled();
    });

    it('clears input after sending', async () => {
        const user = userEvent.setup();
        renderWithTheme(<EllieChatPanel open={true} onClose={vi.fn()} />);
        const input = screen.getByLabelText(/Message Ellie/);
        await user.type(input, 'test message');
        await user.click(screen.getByLabelText('send message'));
        // Input should be cleared after send
        expect(input).toHaveValue('');
    });

    it('has no a11y violations', async () => {
        const { container } = renderWithTheme(<EllieChatPanel open={true} onClose={vi.fn()} />);
        expect(await axe(container)).toHaveNoViolations();
    });
});
