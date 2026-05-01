import { describe, it, expect, vi } from 'vitest';
import { screen } from '@testing-library/react';
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

    it('has no a11y violations', async () => {
        const { container } = renderWithTheme(<HelpPanel open={true} onClose={vi.fn()} />);
        expect(await axe(container)).toHaveNoViolations();
    });
});
