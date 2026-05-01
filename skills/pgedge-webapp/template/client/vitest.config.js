import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

export default defineConfig({
    plugins: [react()],
    test: {
        environment: 'jsdom',
        globals: true,
        setupFiles: ['./src/test/setup.ts'],
        include: ['src/**/*.{test,spec}.{js,jsx,ts,tsx}'],
        coverage: {
            provider: 'v8',
            reporter: ['text', 'json', 'html', 'lcov'],
            exclude: ['node_modules/', 'src/test/', '**/*.config.*', 'src/main.tsx'],
            thresholds: {
                lines: 90, branches: 90, functions: 90, statements: 90,
            },
        },
    },
});
