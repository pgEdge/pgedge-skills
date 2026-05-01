import { useState } from 'react';
import { Box, Card, CardContent, Typography, TextField } from '@mui/material';
import type { Theme } from '@mui/material/styles';
import { SELECT_FIELD_DEFAULT_BG_SX } from './shared/formStyles';
import { useAuth } from '../contexts/useAuth';

const cardSx = { mt: 4, mx: 'auto', maxWidth: 720 };

const titleSx = { mb: 1, fontWeight: 600 };

const bodySx = { color: 'text.secondary', mb: 3 };

const textFieldSx = (theme: Theme) => ({
    '& .MuiOutlinedInput-root': {
        '&:hover .MuiOutlinedInput-notchedOutline': { borderColor: theme.palette.grey[400] },
        '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
            borderColor: theme.palette.primary.main, borderWidth: 2,
        },
    },
    '& .MuiInputLabel-root.Mui-focused': { color: 'primary.main' },
});

const WelcomeCard = (): React.ReactElement => {
    const { user } = useAuth();
    const [demo, setDemo] = useState('');
    return (
        <Box component="main" id="main">
            <Card sx={cardSx}>
                <CardContent>
                    <Typography variant="h4" component="h1" sx={titleSx}>
                        Welcome, {user?.username}
                    </Typography>
                    <Typography variant="body1" sx={bodySx}>
                        This is the bare scaffolding for <PROJECT_NAME>. Replace this
                        component with your application's main content.
                    </Typography>
                    <TextField
                        fullWidth
                        label="Try me"
                        value={demo}
                        onChange={(e) => { setDemo(e.target.value); }}
                        InputLabelProps={{ shrink: true }}
                        sx={(theme) => ({ ...textFieldSx(theme), ...SELECT_FIELD_DEFAULT_BG_SX })}
                    />
                </CardContent>
            </Card>
        </Box>
    );
};

export default WelcomeCard;
