import { describe, it, expect, vi } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../../test/renderWithTheme';
import { axe } from '../../../test/axe';
import SettingsPanel from '../index';

describe('SettingsPanel', () => {
    it('renders the title and the General nav item', () => {
        renderWithTheme(<SettingsPanel open={true} onClose={vi.fn()} />);
        expect(screen.getByRole('dialog', { name: /settings/i })).toBeInTheDocument();
        expect(screen.getByRole('navigation', { name: 'settings sections' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /general/i })).toBeInTheDocument();
    });

    it('calls onClose when close button is clicked', async () => {
        const onClose = vi.fn();
        const user = userEvent.setup();
        renderWithTheme(<SettingsPanel open={true} onClose={onClose} />);
        await user.click(screen.getByLabelText('close settings'));
        expect(onClose).toHaveBeenCalled();
    });

    it('selects the first item on open and renders its content', () => {
        renderWithTheme(<SettingsPanel open={true} onClose={vi.fn()} />);
        // The General placeholder content should be visible in the main pane.
        expect(screen.getByText(/where settings go/i)).toBeInTheDocument();
    });

    it('updates selection when a nav item is clicked', async () => {
        const user = userEvent.setup();
        renderWithTheme(<SettingsPanel open={true} onClose={vi.fn()} />);
        const generalBtn = screen.getByRole('button', { name: /general/i });
        await user.click(generalBtn);
        expect(screen.getByText(/where settings go/i)).toBeInTheDocument();
    });

    it('does not render dialog content when closed', () => {
        renderWithTheme(<SettingsPanel open={false} onClose={vi.fn()} />);
        expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });

    it('has no a11y violations', async () => {
        const { container } = renderWithTheme(<SettingsPanel open={true} onClose={vi.fn()} />);
        expect(await axe(container)).toHaveNoViolations();
    });
});
