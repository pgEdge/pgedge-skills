/*-------------------------------------------------------------------------
 *
 * <PROJECT_NAME>
 *
 * Copyright (c) <CURRENT_YEAR>, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { logger } from '../logger';

describe('logger', () => {
    let errorSpy: ReturnType<typeof vi.spyOn>;
    let warnSpy: ReturnType<typeof vi.spyOn>;
    let infoSpy: ReturnType<typeof vi.spyOn>;
    let debugSpy: ReturnType<typeof vi.spyOn>;

    beforeEach(() => {
        errorSpy = vi.spyOn(console, 'error').mockImplementation(() => undefined);
        warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => undefined);
        infoSpy = vi.spyOn(console, 'info').mockImplementation(() => undefined);
        debugSpy = vi.spyOn(console, 'debug').mockImplementation(() => undefined);
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    it('logger.error calls console.error with all args', () => {
        const err = new Error('test error');
        logger.error('something broke', err);
        expect(errorSpy).toHaveBeenCalledOnce();
        expect(errorSpy).toHaveBeenCalledWith('something broke', err);
    });

    it('logger.warn calls console.warn with all args', () => {
        logger.warn('watch out', { code: 42 });
        expect(warnSpy).toHaveBeenCalledOnce();
        expect(warnSpy).toHaveBeenCalledWith('watch out', { code: 42 });
    });

    it('logger.info calls console.info with all args', () => {
        logger.info('started', 'version', '1.0');
        expect(infoSpy).toHaveBeenCalledOnce();
        expect(infoSpy).toHaveBeenCalledWith('started', 'version', '1.0');
    });

    it('logger.debug calls console.debug with all args', () => {
        logger.debug('debug payload', { key: 'val' });
        expect(debugSpy).toHaveBeenCalledOnce();
        expect(debugSpy).toHaveBeenCalledWith('debug payload', { key: 'val' });
    });

    it('logger.error works with no extra args', () => {
        logger.error('bare message');
        expect(errorSpy).toHaveBeenCalledWith('bare message');
    });

    it('logger.warn works with no extra args', () => {
        logger.warn('bare warning');
        expect(warnSpy).toHaveBeenCalledWith('bare warning');
    });
});
