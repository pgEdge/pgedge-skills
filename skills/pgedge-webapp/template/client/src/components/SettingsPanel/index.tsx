import { useState } from 'react';
import {
    Dialog, DialogTitle, DialogContent, IconButton, Tabs, Tab, Box, Typography,
} from '@mui/material';
import { Close as CloseIcon } from '@mui/icons-material';

interface SettingsPanelProps {
    open: boolean;
    onClose: () => void;
}

const titleSx = { display: 'flex', alignItems: 'center', justifyContent: 'space-between' };
const tabPanelSx = { p: 3 };

const SettingsPanel = ({ open, onClose }: SettingsPanelProps): React.ReactElement => {
    const [tab, setTab] = useState(0);
    return (
        <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth
            aria-labelledby="settings-title">
            <DialogTitle id="settings-title" sx={titleSx}>
                Settings
                <IconButton onClick={onClose} aria-label="close settings">
                    <CloseIcon />
                </IconButton>
            </DialogTitle>
            <DialogContent dividers>
                <Tabs value={tab} onChange={(_, v) => { setTab(v); }}
                    aria-label="settings tabs">
                    <Tab label="General" />
                </Tabs>
                <Box role="tabpanel" sx={tabPanelSx}>
                    <Typography variant="body1">
                        This is where settings go. Add tabs and content as your application
                        grows.
                    </Typography>
                </Box>
            </DialogContent>
        </Dialog>
    );
};

export default SettingsPanel;
