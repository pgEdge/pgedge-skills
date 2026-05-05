import { useState } from 'react';
import { Drawer, Box, Typography, IconButton, Divider } from '@mui/material';
import { Close as CloseIcon } from '@mui/icons-material';
import ChatInput from './ChatInput';
import ChatMessage, { type Message } from './ChatMessage';

interface EllieChatPanelProps {
    open: boolean;
    onClose: () => void;
}

const drawerPaperSx = { width: { xs: '100vw', sm: 480 }, maxWidth: '100vw' };
const headerSx = { display: 'flex', alignItems: 'center', justifyContent: 'space-between', p: 2 };
const messagesSx = { flex: 1, overflowY: 'auto', p: 2 };
const containerSx = { height: '100%', display: 'flex', flexDirection: 'column' };

const EllieChatPanel = ({ open, onClose }: EllieChatPanelProps): React.ReactElement => {
    const [messages, setMessages] = useState<Message[]>([]);
    const handleSend = (content: string) => {
        // TODO: wire to your chat backend.
        setMessages((m) => [...m, {
            id: String(Date.now()), role: 'user', content,
        }]);
    };
    return (
        <Drawer anchor="right" open={open} onClose={onClose}
            PaperProps={{ sx: drawerPaperSx }}
            aria-labelledby="ellie-title">
            <Box sx={containerSx}>
                <Box sx={headerSx}>
                    <Typography id="ellie-title" variant="h6">Ellie</Typography>
                    <IconButton onClick={onClose} aria-label="close chat"><CloseIcon /></IconButton>
                </Box>
                <Divider />
                <Box sx={messagesSx} role="log" aria-live="polite">
                    {messages.map((m) => <ChatMessage key={m.id} message={m} />)}
                </Box>
                <Divider />
                <ChatInput onSend={handleSend} />
            </Box>
        </Drawer>
    );
};

export default EllieChatPanel;
