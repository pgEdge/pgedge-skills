/*-------------------------------------------------------------------------
 *
 * <PROJECT_NAME>
 *
 * Copyright (c) <CURRENT_YEAR>, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

/**
 * Centralized logging utility.
 *
 * Every production module must use this logger instead of calling
 * console methods directly.  ESLint enforces this via the
 * no-console rule; only this file is exempt.
 */
export const logger = {
    error(...args: unknown[]): void {
        console.error(...args);
    },
    warn(...args: unknown[]): void {
        console.warn(...args);
    },
    info(...args: unknown[]): void {
        console.info(...args);
    },
    debug(...args: unknown[]): void {
        console.debug(...args);
    },
};
