import { describe, it, expect, vi, beforeEach } from 'vitest';
import { waitFor } from '@testing-library/react';
import { renderWithTheme } from '../../test/renderWithTheme';
import { axe } from '../../test/axe';
import App from '../../App';

describe('App a11y sweep', () => {
    beforeEach(() => { global.fetch = vi.fn() as unknown as typeof fetch; });

    it('login state has no a11y violations', async () => {
        (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
            ok: false, json: () => Promise.resolve({ authenticated: false }),
        });
        const { container } = renderWithTheme(<App />);
        await waitFor(() => container.querySelector('h5'));
        expect(await axe(container)).toHaveNoViolations();
    });

    it('authenticated state has no a11y violations', async () => {
        (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
            ok: true, json: () => Promise.resolve({ authenticated: true, username: 'tester' }),
        });
        const { container } = renderWithTheme(<App />);
        await waitFor(() => container.querySelector('h4'));
        expect(await axe(container)).toHaveNoViolations();
    });
});
