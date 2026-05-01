import { axe } from 'vitest-axe';
import * as matchers from 'vitest-axe/matchers';
import { expect } from 'vitest';

expect.extend(matchers);

declare module '@vitest/expect' {
    interface Assertion<T> {
        toHaveNoViolations(): T;
    }
}

export { axe };
