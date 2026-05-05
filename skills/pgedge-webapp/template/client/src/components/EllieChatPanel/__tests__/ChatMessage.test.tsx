import { describe, it, expect } from 'vitest';
import { renderWithTheme } from '../../../test/renderWithTheme';
import ChatMessage from '../ChatMessage';
import type { Message } from '../ChatMessage';
import { screen } from '@testing-library/react';

describe('ChatMessage', () => {
    it('renders a user message', () => {
        const msg: Message = { id: '1', role: 'user', content: 'Hello from user' };
        renderWithTheme(<ChatMessage message={msg} />);
        expect(screen.getByText('Hello from user')).toBeInTheDocument();
    });

    it('renders an assistant message', () => {
        const msg: Message = { id: '2', role: 'assistant', content: 'Hello from assistant' };
        renderWithTheme(<ChatMessage message={msg} />);
        expect(screen.getByText('Hello from assistant')).toBeInTheDocument();
    });
});
