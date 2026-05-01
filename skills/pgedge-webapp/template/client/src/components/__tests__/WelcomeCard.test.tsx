import { describe, it, expect, vi } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../test/renderWithTheme';
import { axe } from '../../test/axe';
import WelcomeCard from '../WelcomeCard';

vi.mock('../../contexts/useAuth', () => ({
    useAuth: () => ({ user: { authenticated: true, username: 'tester' }, loading: false }),
}));

describe('WelcomeCard', () => {
    it('greets the authenticated user', () => {
        renderWithTheme(<WelcomeCard />);
        expect(screen.getByRole('heading', { name: /Welcome, tester/i })).toBeInTheDocument();
    });

    it('accepts text in the demo field', async () => {
        const user = userEvent.setup();
        renderWithTheme(<WelcomeCard />);
        const input = screen.getByLabelText(/Try me/i);
        await user.type(input, 'hello');
        expect(input).toHaveValue('hello');
    });

    it('has no accessibility violations', async () => {
        const { container } = renderWithTheme(<WelcomeCard />);
        expect(await axe(container)).toHaveNoViolations();
    });
});
