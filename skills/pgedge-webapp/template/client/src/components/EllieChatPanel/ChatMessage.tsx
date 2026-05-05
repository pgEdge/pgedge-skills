import { Box, Paper, Typography } from '@mui/material';
import type { Theme } from '@mui/material/styles';

export interface Message {
    id: string;
    role: 'user' | 'assistant';
    content: string;
}

interface ChatMessageProps { message: Message }

const rowSx = (role: Message['role']) => ({
    display: 'flex',
    justifyContent: role === 'user' ? 'flex-end' : 'flex-start',
    mb: 1.5,
});
const bubbleSx = (theme: Theme, role: Message['role']) => ({
    maxWidth: '80%', p: 1.5, borderRadius: 2,
    bgcolor: role === 'user' ? theme.palette.primary.main : theme.palette.background.paper,
    color: role === 'user' ? theme.palette.primary.contrastText : theme.palette.text.primary,
    border: role === 'assistant' ? `1px solid ${theme.palette.divider}` : 'none',
});

const ChatMessage = ({ message }: ChatMessageProps): React.ReactElement => (
    <Box sx={rowSx(message.role)}>
        <Paper elevation={0} sx={(theme) => bubbleSx(theme, message.role)}>
            <Typography variant="body2">{message.content}</Typography>
        </Paper>
    </Box>
);

export default ChatMessage;
