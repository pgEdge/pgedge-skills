import { axe, toHaveNoViolations } from 'vitest-axe';
import { expect } from 'vitest';

expect.extend({ toHaveNoViolations });

export { axe };
