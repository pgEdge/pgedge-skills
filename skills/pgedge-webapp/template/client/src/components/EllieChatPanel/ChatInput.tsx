import { useState } from 'react';
import { Box, TextField, IconButton } from '@mui/material';
import { Send as SendIcon } from '@mui/icons-material';
import type { Theme } from '@mui/material/styles';
import { SELECT_FIELD_SX } from '../shared/formStyles';

interface ChatInputProps {
    onSend: (text: string) => void;
}

const containerSx = { display: 'flex', alignItems: 'flex-end', gap: 1, p: 2 };

const textFieldSx = (theme: Theme) => ({
    flex: 1,
    '& .MuiOutlinedInput-root': {
        '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
            borderColor: theme.palette.primary.main, borderWidth: 2,
        },
    },
});

const ChatInput = ({ onSend }: ChatInputProps): React.ReactElement => {
    const [text, setText] = useState('');
    const submit = () => {
        if (!text.trim()) { return; }
        onSend(text);
        setText('');
    };
    return (
        <Box sx={containerSx}>
            <TextField
                multiline maxRows={4}
                label="Message Ellie"
                value={text}
                onChange={(e) => { setText(e.target.value); }}
                InputLabelProps={{ shrink: true }}
                sx={(theme) => ({ ...textFieldSx(theme), ...SELECT_FIELD_SX })}
            />
            <IconButton color="primary" aria-label="send message"
                onClick={submit} disabled={!text.trim()}>
                <SendIcon />
            </IconButton>
        </Box>
    );
};

export default ChatInput;
