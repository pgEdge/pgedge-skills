import { useState, useEffect, useMemo } from 'react';
import { ThemeProvider } from '@mui/material/styles';
import { Box, CircularProgress, CssBaseline, type PaletteMode } from '@mui/material';
import { AuthProvider } from './contexts/AuthContext';
import { useAuth } from './contexts/useAuth';
import ErrorBoundary from './components/ErrorBoundary';
import Header from './components/Header';
import Login from './components/Login';
import WelcomeCard from './components/WelcomeCard';
import { createPgedgeTheme, loginTheme } from './theme/pgedgeTheme';

const styles = {
    loadingContainer: {
        minHeight: '100vh',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        bgcolor: 'background.default',
    },
    mainLayoutRoot: {
        height: '100vh', display: 'flex', flexDirection: 'column' as const,
        bgcolor: 'background.default', overflow: 'hidden',
    },
    mainLayoutBody: { flex: 1, overflow: 'auto', position: 'relative' as const },
    skipLink: {
        position: 'absolute' as const, left: -9999, top: 8,
        '&:focus': { left: 8, zIndex: 9999 },
        bgcolor: 'background.paper', p: 1, borderRadius: 1,
    },
};

const AppContent = () => {
    const [mode, setMode] = useState<PaletteMode>(() => {
        const saved = localStorage.getItem('theme-mode');
        return saved === 'dark' ? 'dark' : 'light';
    });
    const { user, loading } = useAuth();

    useEffect(() => { localStorage.setItem('theme-mode', mode); }, [mode]);

    const theme = useMemo(() => createPgedgeTheme(mode), [mode]);

    if (loading) {
        return (
            <ThemeProvider theme={theme}>
                <CssBaseline />
                <Box sx={styles.loadingContainer}>
                    <CircularProgress aria-label="Loading application" />
                </Box>
            </ThemeProvider>
        );
    }

    if (!user) {
        return (
            <ThemeProvider theme={loginTheme}>
                <CssBaseline />
                <Login />
            </ThemeProvider>
        );
    }

    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <ErrorBoundary>
                <Box sx={styles.mainLayoutRoot}>
                    <Box component="a" href="#main" sx={styles.skipLink}>
                        Skip to main content
                    </Box>
                    <Header onToggleTheme={() => setMode((m) => (m === 'light' ? 'dark' : 'light'))} />
                    <Box sx={styles.mainLayoutBody}>
                        <WelcomeCard />
                    </Box>
                </Box>
            </ErrorBoundary>
        </ThemeProvider>
    );
};

const App = (): React.ReactElement => (
    <ErrorBoundary>
        <AuthProvider>
            <AppContent />
        </AuthProvider>
    </ErrorBoundary>
);

export default App;
