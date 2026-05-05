import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithTheme } from '../../test/renderWithTheme';
import ErrorBoundary from '../ErrorBoundary';

// A component that throws when shouldThrow is true
const ThrowingChild = ({ shouldThrow }: { shouldThrow: boolean }): React.ReactElement => {
    if (shouldThrow) {
        throw new Error('Test render error');
    }
    return <div>no error</div>;
};

// A component that always throws for fallback tests
const AlwaysThrows = (): React.ReactElement => {
    throw new Error('Always fails');
};

describe('ErrorBoundary', () => {
    let consoleErrorSpy: ReturnType<typeof vi.spyOn>;

    beforeEach(() => {
        // Suppress React's error boundary console output
        consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    });

    afterEach(() => {
        consoleErrorSpy.mockRestore();
    });

    it('renders children when there is no error', () => {
        renderWithTheme(
            <ErrorBoundary>
                <div>hello</div>
            </ErrorBoundary>,
        );
        expect(screen.getByText('hello')).toBeInTheDocument();
    });

    it('renders fallback UI when a child throws', () => {
        renderWithTheme(
            <ErrorBoundary>
                <ThrowingChild shouldThrow={true} />
            </ErrorBoundary>,
        );
        expect(screen.getByText(/Something went wrong/i)).toBeInTheDocument();
        expect(screen.getByText(/The application has crashed/i)).toBeInTheDocument();
    });

    it('renders custom fallback when provided', () => {
        renderWithTheme(
            <ErrorBoundary fallback={<div>custom fallback</div>}>
                <AlwaysThrows />
            </ErrorBoundary>,
        );
        expect(screen.getByText('custom fallback')).toBeInTheDocument();
    });

    it('shows error details by default', () => {
        renderWithTheme(
            <ErrorBoundary>
                <AlwaysThrows />
            </ErrorBoundary>,
        );
        expect(screen.getByTestId('error-boundary-details')).toBeInTheDocument();
        expect(screen.getByTestId('error-boundary-details').textContent).toContain('Always fails');
    });

    it('toggles error details when clicking the details button', async () => {
        const user = userEvent.setup();
        renderWithTheme(
            <ErrorBoundary>
                <AlwaysThrows />
            </ErrorBoundary>,
        );
        // Initially open (detailsOpen: true), button label says "Hide error details"
        const hideBtn = screen.getByLabelText('Hide error details');
        expect(hideBtn).toBeInTheDocument();

        await user.click(hideBtn);
        // After clicking, it should say "Show error details"
        expect(screen.getByLabelText('Show error details')).toBeInTheDocument();

        // Click again to re-open
        await user.click(screen.getByLabelText('Show error details'));
        expect(screen.getByLabelText('Hide error details')).toBeInTheDocument();
    });

    it('shows reload button', () => {
        renderWithTheme(
            <ErrorBoundary>
                <AlwaysThrows />
            </ErrorBoundary>,
        );
        expect(screen.getByRole('button', { name: /reload/i })).toBeInTheDocument();
    });

    it('calls window.location.reload when Reload is clicked', async () => {
        const user = userEvent.setup();
        // jsdom does not allow spying on location.reload directly;
        // replace the whole location object instead.
        const origLocation = window.location;
        const reloadMock = vi.fn();
        Object.defineProperty(window, 'location', {
            value: { ...origLocation, reload: reloadMock },
            writable: true,
            configurable: true,
        });

        renderWithTheme(
            <ErrorBoundary>
                <AlwaysThrows />
            </ErrorBoundary>,
        );
        await user.click(screen.getByRole('button', { name: /reload/i }));
        expect(reloadMock).toHaveBeenCalledOnce();

        // Restore
        Object.defineProperty(window, 'location', {
            value: origLocation,
            writable: true,
            configurable: true,
        });
    });

    it('calls onError callback when a child throws', () => {
        const onError = vi.fn();
        renderWithTheme(
            <ErrorBoundary onError={onError}>
                <AlwaysThrows />
            </ErrorBoundary>,
        );
        expect(onError).toHaveBeenCalledOnce();
        expect(onError.mock.calls[0][0]).toBeInstanceOf(Error);
        expect(onError.mock.calls[0][0].message).toBe('Always fails');
    });

    it('swallows exceptions thrown by onError callback', () => {
        const throwingCallback = vi.fn().mockImplementation(() => {
            throw new Error('callback broke');
        });
        // Should not throw to the test
        expect(() => {
            renderWithTheme(
                <ErrorBoundary onError={throwingCallback}>
                    <AlwaysThrows />
                </ErrorBoundary>,
            );
        }).not.toThrow();
    });

    it('shows error text in details even without a stack', () => {
        // Error with no stack property
        const NoStackChild = (): React.ReactElement => {
            const err = new Error('No stack error');
            err.stack = undefined as unknown as string;
            throw err;
        };
        renderWithTheme(
            <ErrorBoundary>
                <NoStackChild />
            </ErrorBoundary>,
        );
        const details = screen.getByTestId('error-boundary-details');
        expect(details.textContent).toContain('No stack error');
    });
});
