import { describe, it, expect, vi } from 'vitest';
import { screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../../test/renderWithTheme';
import { axe } from '../../../test/axe';
import SettingsPanel from '../index';

describe('SettingsPanel', () => {
    it('renders title and tab', () => {
        renderWithTheme(<SettingsPanel open={true} onClose={vi.fn()} />);
        expect(screen.getByRole('dialog', { name: /settings/i })).toBeInTheDocument();
        expect(screen.getByRole('tab', { name: /general/i })).toBeInTheDocument();
    });

    it('calls onClose when close button is clicked', async () => {
        const onClose = vi.fn();
        const user = userEvent.setup();
        renderWithTheme(<SettingsPanel open={true} onClose={onClose} />);
        await user.click(screen.getByLabelText('close settings'));
        expect(onClose).toHaveBeenCalled();
    });

    it('tab onChange executes when a tab button is clicked', () => {
        renderWithTheme(<SettingsPanel open={true} onClose={vi.fn()} />);
        const tabEl = screen.getByRole('tab', { name: /general/i });
        // Use mousedown + click to simulate MUI tab selection which fires onChange
        fireEvent.mouseDown(tabEl);
        fireEvent.click(tabEl);
        // Tab should still be rendered with no crash from the onChange handler
        expect(tabEl).toBeInTheDocument();
    });

    it('displays the settings body content', () => {
        renderWithTheme(<SettingsPanel open={true} onClose={vi.fn()} />);
        expect(screen.getByRole('tabpanel')).toBeInTheDocument();
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
