/*-------------------------------------------------------------------------
 *
 * <PROJECT_NAME>
 *
 * Copyright (c) <CURRENT_YEAR>, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

import { describe, it, expect } from 'vitest';
import { createPgedgeTheme, loginTheme } from '../pgedgeTheme';

describe('createPgedgeTheme', () => {
    it('defaults to light mode', () => {
        const theme = createPgedgeTheme();
        expect(theme.palette.mode).toBe('light');
    });

    it('creates a light theme when mode is "light"', () => {
        const theme = createPgedgeTheme('light');
        expect(theme.palette.mode).toBe('light');
        expect(theme.palette.background.default).toBe('#F9FAFB');
        expect(theme.palette.background.paper).toBe('#FFFFFF');
        expect(theme.palette.text.primary).toBe('#1F2937');
    });

    it('creates a dark theme when mode is "dark"', () => {
        const theme = createPgedgeTheme('dark');
        expect(theme.palette.mode).toBe('dark');
        expect(theme.palette.background.default).toBe('#0F172A');
        expect(theme.palette.background.paper).toBe('#1E293B');
        expect(theme.palette.text.primary).toBe('#F1F5F9');
    });

    it('dark theme has brightened primary color', () => {
        const darkTheme = createPgedgeTheme('dark');
        const lightTheme = createPgedgeTheme('light');
        // dark mode uses cyan.5 (#22B8CF) instead of cyan.6 (#15AABF)
        expect(darkTheme.palette.primary.main).toBe('#22B8CF');
        expect(lightTheme.palette.primary.main).toBe('#15AABF');
    });

    it('dark theme custom accent differs from light', () => {
        const darkTheme = createPgedgeTheme('dark');
        const lightTheme = createPgedgeTheme('light');
        expect(darkTheme.palette.custom.accent).toBe('#22B8CF');
        expect(lightTheme.palette.custom.accent).toBe('#15AABF');
    });

    it('dark theme custom status colors differ from light', () => {
        const darkTheme = createPgedgeTheme('dark');
        const lightTheme = createPgedgeTheme('light');
        expect(darkTheme.palette.custom.status.connected).toBe('#10B981');
        expect(lightTheme.palette.custom.status.connected).toBe('#059669');
    });

    it('theme has expected shape border radius', () => {
        const theme = createPgedgeTheme('light');
        expect(theme.shape.borderRadius).toBe(8);
    });

    it('theme has Inter as primary font', () => {
        const theme = createPgedgeTheme('light');
        expect(theme.typography.fontFamily).toContain('Inter');
    });

    it('loginTheme is a light theme', () => {
        expect(loginTheme.palette.mode).toBe('light');
    });

    it('both modes have MuiButton override', () => {
        const lightTheme = createPgedgeTheme('light');
        const darkTheme = createPgedgeTheme('dark');
        expect(lightTheme.components?.MuiButton).toBeDefined();
        expect(darkTheme.components?.MuiButton).toBeDefined();
    });

    it('dark theme divider color differs from light', () => {
        const darkTheme = createPgedgeTheme('dark');
        const lightTheme = createPgedgeTheme('light');
        expect(darkTheme.palette.divider).toBe('#334155');
        expect(lightTheme.palette.divider).toBe('#D1D5DB');
    });
});
