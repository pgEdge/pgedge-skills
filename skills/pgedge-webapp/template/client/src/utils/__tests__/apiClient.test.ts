/*-------------------------------------------------------------------------
 *
 * <PROJECT_NAME>
 *
 * Copyright (c) <CURRENT_YEAR>, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
    ApiError,
    apiGet,
    apiPost,
    apiPut,
    apiPatch,
    apiDelete,
    apiFetch,
    onDisconnect,
    resetConnectionHealth,
} from '../apiClient';

function makeFetchResponse(
    options: {
        ok?: boolean;
        status?: number;
        json?: unknown;
        text?: string;
    } = {},
): Response {
    const { ok = true, status = 200, json, text = '' } = options;
    return {
        ok,
        status,
        json: () => Promise.resolve(json),
        text: () => Promise.resolve(text),
    } as unknown as Response;
}

describe('ApiError', () => {
    it('constructs with statusCode and errorBody', () => {
        const err = new ApiError('bad', 404, 'not found');
        expect(err.name).toBe('ApiError');
        expect(err.statusCode).toBe(404);
        expect(err.errorBody).toBe('not found');
        expect(err.message).toBe('bad');
    });

    it('has empty errorBody by default', () => {
        const err = new ApiError('oops', 500);
        expect(err.errorBody).toBe('');
    });

    it('is instanceof Error', () => {
        expect(new ApiError('x', 400)).toBeInstanceOf(Error);
    });
});

describe('onDisconnect / resetConnectionHealth', () => {
    beforeEach(() => {
        resetConnectionHealth();
        globalThis.fetch = vi.fn();
    });

    it('returns an unsubscribe function', () => {
        const listener = vi.fn();
        const unsub = onDisconnect(listener);
        expect(typeof unsub).toBe('function');
        unsub();
    });

    it('fires listener on auth failure (401)', async () => {
        const listener = vi.fn();
        onDisconnect(listener);
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 401, text: '{"error":"Unauthorized"}' }),
        );
        await expect(apiGet('/test')).rejects.toThrow();
        expect(listener).toHaveBeenCalledWith('auth');
    });

    it('fires listener on 3 consecutive server failures', async () => {
        resetConnectionHealth();
        const listener = vi.fn();
        onDisconnect(listener);
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 500, text: '' }),
        );
        // 3 failures needed to trigger
        await expect(apiGet('/test')).rejects.toThrow();
        await expect(apiGet('/test')).rejects.toThrow();
        expect(listener).not.toHaveBeenCalled();
        await expect(apiGet('/test')).rejects.toThrow();
        expect(listener).toHaveBeenCalledWith('server');
    });

    it('does not fire listener twice once disconnected', async () => {
        resetConnectionHealth();
        const listener = vi.fn();
        onDisconnect(listener);
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 401, text: '' }),
        );
        await expect(apiGet('/test')).rejects.toThrow();
        await expect(apiGet('/test')).rejects.toThrow();
        // listener should only fire once (first auth failure triggers it)
        expect(listener).toHaveBeenCalledTimes(1);
    });

    it('unsubscribe prevents listener from being called', async () => {
        resetConnectionHealth();
        const listener = vi.fn();
        const unsub = onDisconnect(listener);
        unsub();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 401, text: '' }),
        );
        await expect(apiGet('/test')).rejects.toThrow();
        expect(listener).not.toHaveBeenCalled();
    });

    it('resetConnectionHealth clears disconnect state', async () => {
        resetConnectionHealth();
        const listener = vi.fn();
        onDisconnect(listener);
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 401, text: '' }),
        );
        await expect(apiGet('/test')).rejects.toThrow();
        expect(listener).toHaveBeenCalledTimes(1);

        // After reset, auth failure should fire again
        resetConnectionHealth();
        await expect(apiGet('/test')).rejects.toThrow();
        expect(listener).toHaveBeenCalledTimes(2);
    });
});

describe('apiGet', () => {
    beforeEach(() => {
        resetConnectionHealth();
        globalThis.fetch = vi.fn();
    });

    it('returns parsed JSON on success', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ json: { data: 42 } }),
        );
        const result = await apiGet<{ data: number }>('/api/v1/things');
        expect(result).toEqual({ data: 42 });
    });

    it('sends GET with credentials include', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ json: {} }),
        );
        await apiGet('/api/v1/things');
        expect(globalThis.fetch).toHaveBeenCalledWith(
            '/api/v1/things',
            expect.objectContaining({ method: 'GET', credentials: 'include' }),
        );
    });

    it('returns raw text when rawResponse is true', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ text: 'hello world' }),
        );
        const result = await apiGet<string>('/api/v1/text', { rawResponse: true });
        expect(result).toBe('hello world');
    });

    it('returns undefined for 204 No Content', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ status: 204 }),
        );
        const result = await apiGet('/api/v1/empty');
        expect(result).toBeUndefined();
    });

    it('throws ApiError with error field from JSON body', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 400, text: '{"error":"Bad Request"}' }),
        );
        await expect(apiGet('/api/v1/bad')).rejects.toMatchObject({
            name: 'ApiError',
            statusCode: 400,
            message: 'Bad Request',
        });
    });

    it('throws ApiError with message field from JSON body', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 422, text: '{"message":"Invalid input"}' }),
        );
        await expect(apiGet('/api/v1/bad')).rejects.toMatchObject({
            name: 'ApiError',
            statusCode: 422,
            message: 'Invalid input',
        });
    });

    it('throws ApiError using raw text when JSON parse fails', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 503, text: 'Service Unavailable' }),
        );
        await expect(apiGet('/api/v1/bad')).rejects.toMatchObject({
            name: 'ApiError',
            statusCode: 503,
            message: 'Service Unavailable',
        });
    });

    it('throws ApiError with fallback message when text is empty', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 503, text: '' }),
        );
        await expect(apiGet('/api/v1/bad')).rejects.toMatchObject({
            name: 'ApiError',
            statusCode: 503,
            message: expect.stringContaining('Request failed'),
        });
    });

    it('throws network errors directly', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockRejectedValue(
            new TypeError('Failed to fetch'),
        );
        await expect(apiGet('/api/v1/things')).rejects.toThrow('Failed to fetch');
    });

    it('passes AbortSignal to fetch', async () => {
        const controller = new AbortController();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ json: {} }),
        );
        await apiGet('/api/v1/things', { signal: controller.signal });
        expect(globalThis.fetch).toHaveBeenCalledWith(
            '/api/v1/things',
            expect.objectContaining({ signal: controller.signal }),
        );
    });

    it('merges extra headers', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ json: {} }),
        );
        await apiGet('/api/v1/things', { headers: { 'X-Custom': 'yes' } });
        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1];
        expect(callArgs.headers).toMatchObject({ 'X-Custom': 'yes' });
    });
});

describe('apiPost', () => {
    beforeEach(() => {
        resetConnectionHealth();
        globalThis.fetch = vi.fn();
    });

    it('sends POST with JSON body', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ json: { id: 1 } }),
        );
        const result = await apiPost<{ id: number }>('/api/v1/items', { name: 'test' });
        expect(result).toEqual({ id: 1 });
        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1];
        expect(callArgs.method).toBe('POST');
        expect(callArgs.headers['Content-Type']).toBe('application/json');
        expect(callArgs.body).toBe('{"name":"test"}');
    });

    it('sends POST without body when omitted', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ json: {} }),
        );
        await apiPost('/api/v1/items');
        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1];
        expect(callArgs.body).toBeUndefined();
        expect(callArgs.headers['Content-Type']).toBeUndefined();
    });

    it('throws ApiError on 4xx', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 403, text: '{"error":"Forbidden"}' }),
        );
        await expect(apiPost('/api/v1/items', {})).rejects.toMatchObject({
            name: 'ApiError',
            statusCode: 403,
        });
    });
});

describe('apiPut', () => {
    beforeEach(() => {
        resetConnectionHealth();
        globalThis.fetch = vi.fn();
    });

    it('sends PUT with body', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ json: { updated: true } }),
        );
        const result = await apiPut<{ updated: boolean }>('/api/v1/items/1', { val: 99 });
        expect(result).toEqual({ updated: true });
        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1];
        expect(callArgs.method).toBe('PUT');
    });

    it('throws on 5xx', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 500, text: '' }),
        );
        await expect(apiPut('/api/v1/items/1', {})).rejects.toMatchObject({ name: 'ApiError' });
    });
});

describe('apiPatch', () => {
    beforeEach(() => {
        resetConnectionHealth();
        globalThis.fetch = vi.fn();
    });

    it('sends PATCH with body', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ json: { patched: true } }),
        );
        const result = await apiPatch<{ patched: boolean }>('/api/v1/items/1', { x: 1 });
        expect(result).toEqual({ patched: true });
        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1];
        expect(callArgs.method).toBe('PATCH');
    });
});

describe('apiDelete', () => {
    beforeEach(() => {
        resetConnectionHealth();
        globalThis.fetch = vi.fn();
    });

    it('sends DELETE without body', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ status: 204 }),
        );
        const result = await apiDelete('/api/v1/items/1');
        expect(result).toBeUndefined();
        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1];
        expect(callArgs.method).toBe('DELETE');
        expect(callArgs.body).toBeUndefined();
    });

    it('throws on 4xx', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 404, text: '{"error":"Not Found"}' }),
        );
        await expect(apiDelete('/api/v1/items/99')).rejects.toMatchObject({
            name: 'ApiError',
            statusCode: 404,
        });
    });
});

describe('apiFetch', () => {
    beforeEach(() => {
        resetConnectionHealth();
        globalThis.fetch = vi.fn();
    });

    it('returns raw Response on success', async () => {
        const fakeResponse = makeFetchResponse({ json: { ok: true } });
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(fakeResponse);
        const result = await apiFetch('/api/v1/stream');
        expect(result).toBe(fakeResponse);
    });

    it('always includes credentials: include', async () => {
        const fakeResponse = makeFetchResponse();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(fakeResponse);
        await apiFetch('/api/v1/stream', { method: 'POST' });
        expect(globalThis.fetch).toHaveBeenCalledWith(
            '/api/v1/stream',
            expect.objectContaining({ credentials: 'include', method: 'POST' }),
        );
    });

    it('records auth failure on 401', async () => {
        const listener = vi.fn();
        onDisconnect(listener);
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 401 }),
        );
        await apiFetch('/api/v1/stream');
        expect(listener).toHaveBeenCalledWith('auth');
    });

    it('records server failure on 500+', async () => {
        resetConnectionHealth();
        const listener = vi.fn();
        onDisconnect(listener);
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
            makeFetchResponse({ ok: false, status: 503 }),
        );
        // 3 failures needed
        await apiFetch('/api/v1/stream');
        await apiFetch('/api/v1/stream');
        await apiFetch('/api/v1/stream');
        expect(listener).toHaveBeenCalledWith('server');
    });

    it('throws on network error', async () => {
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockRejectedValue(
            new TypeError('Network failure'),
        );
        await expect(apiFetch('/api/v1/stream')).rejects.toThrow('Network failure');
    });

    it('merges caller init with credentials', async () => {
        const fakeResponse = makeFetchResponse();
        (globalThis.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(fakeResponse);
        await apiFetch('/api/v1/stream', { headers: { Authorization: 'Bearer tok' } });
        const callArgs = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1];
        expect(callArgs.headers).toMatchObject({ Authorization: 'Bearer tok' });
    });
});
