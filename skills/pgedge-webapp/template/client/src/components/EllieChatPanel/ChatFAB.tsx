import { Fab, Tooltip } from '@mui/material';
import { Chat as ChatIcon } from '@mui/icons-material';

interface ChatFABProps { onClick: () => void; isOpen: boolean }

const fabSx = { position: 'fixed' as const, bottom: 24, right: 24 };

const ChatFAB = ({ onClick, isOpen }: ChatFABProps): React.ReactElement | null => {
    if (isOpen) { return null; }
    return (
        <Tooltip title="Open Ellie">
            <Fab color="primary" aria-label="open chat" onClick={onClick} sx={fabSx}>
                <ChatIcon />
            </Fab>
        </Tooltip>
    );
};

export default ChatFAB;
