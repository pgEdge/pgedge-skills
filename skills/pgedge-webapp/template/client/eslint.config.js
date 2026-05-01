import js from '@eslint/js';
import globals from 'globals';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import jsxA11y from 'eslint-plugin-jsx-a11y';
import tseslint from 'typescript-eslint';

export default tseslint.config(
    { ignores: ['dist', 'coverage', 'node_modules'] },
    {
        extends: [js.configs.recommended, ...tseslint.configs.recommended],
        files: ['**/*.{ts,tsx}'],
        languageOptions: {
            ecmaVersion: 2020,
            globals: globals.browser,
        },
        plugins: {
            'react-hooks': reactHooks,
            'react-refresh': reactRefresh,
            'jsx-a11y': jsxA11y,
        },
        rules: {
            ...reactHooks.configs.recommended.rules,
            ...jsxA11y.configs.recommended.rules,
            'react-hooks/rules-of-hooks': 'error',
            'react-hooks/exhaustive-deps': 'error',
            'react-refresh/only-export-components': ['warn', { allowConstantExport: true }],
            '@typescript-eslint/no-unused-vars': ['error', {
                argsIgnorePattern: '^_', varsIgnorePattern: '^_',
            }],
            '@typescript-eslint/no-explicit-any': 'error',
            'no-console': 'error',
            'prefer-const': 'error',
            'no-var': 'error',
            eqeqeq: ['error', 'always', { null: 'ignore' }],
            curly: ['error', 'all'],
            'jsx-a11y/alt-text': 'error',
            'jsx-a11y/label-has-associated-control': 'error',
            'jsx-a11y/no-noninteractive-element-interactions': 'error',
            'jsx-a11y/click-events-have-key-events': 'error',
        },
    },
    {
        files: ['src/utils/logger.ts'],
        rules: { 'no-console': 'off' },
    }
);
