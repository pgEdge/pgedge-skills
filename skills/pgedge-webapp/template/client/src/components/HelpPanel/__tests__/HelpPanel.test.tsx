import { describe, it, expect, vi } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../../test/renderWithTheme';
import { axe } from '../../../test/axe';
import HelpPanel from '../index';

describe('HelpPanel', () => {
    it('renders the overview page when open', () => {
        renderWithTheme(<HelpPanel open={true} onClose={vi.fn()} />);
        expect(screen.getByRole('heading', { name: /Overview/i })).toBeInTheDocument();
    });

    it('exposes a close button', () => {
        renderWithTheme(<HelpPanel open={true} onClose={vi.fn()} />);
        expect(screen.getByLabelText('close help')).toBeInTheDocument();
    });

    it('calls onClose when close button clicked', async () => {
        const onClose = vi.fn();
        const user = userEvent.setup();
        renderWithTheme(<HelpPanel open={true} onClose={onClose} />);
        await user.click(screen.getByLabelText('close help'));
        expect(onClose).toHaveBeenCalled();
    });

    it('shows the nav sidebar with page list items', () => {
        renderWithTheme(<HelpPanel open={true} onClose={vi.fn()} />);
        // The nav should contain the Overview page link
        expect(screen.getByRole('navigation', { name: 'help pages' })).toBeInTheDocument();
        // "Overview" appears in both the nav and as a heading - use getAllByText
        expect(screen.getAllByText('Overview').length).toBeGreaterThanOrEqual(1);
    });

    it('clicking a nav item keeps it selected', async () => {
        const user = userEvent.setup();
        renderWithTheme(<HelpPanel open={true} onClose={vi.fn()} />);
        const overviewBtn = screen.getByRole('button', { name: 'Overview' });
        await user.click(overviewBtn);
        // Still shows overview content after clicking the already-active page
        expect(screen.getByRole('heading', { name: /Overview/i })).toBeInTheDocument();
    });

    it('has no a11y violations', async () => {
        const { container } = renderWithTheme(<HelpPanel open={true} onClose={vi.fn()} />);
        expect(await axe(container)).toHaveNoViolations();
    });
});
