import { useState } from 'react';
import {
    Drawer, Box, List, ListItemButton, ListItemText, Typography, IconButton, Divider,
} from '@mui/material';
import { Close as CloseIcon } from '@mui/icons-material';
import { pages } from './pages';

interface HelpPanelProps {
    open: boolean;
    onClose: () => void;
}

const drawerPaperSx = { width: { xs: '100vw', sm: 600 }, maxWidth: '100vw' };
const headerSx = { display: 'flex', alignItems: 'center', justifyContent: 'space-between', p: 2 };
const bodySx = { display: 'flex', height: 'calc(100% - 64px)' };
const railSx = { width: 200, borderRight: '1px solid', borderColor: 'divider', overflowY: 'auto' };
const contentSx = { flex: 1, p: 3, overflowY: 'auto' };

const HelpPanel = ({ open, onClose }: HelpPanelProps): React.ReactElement => {
    const [activeKey, setActiveKey] = useState(pages[0]?.key ?? '');
    const Active = pages.find((p) => p.key === activeKey)?.component;
    return (
        <Drawer anchor="right" open={open} onClose={onClose}
            PaperProps={{ sx: drawerPaperSx }}
            aria-labelledby="help-title">
            <Box sx={headerSx}>
                <Typography variant="h6" id="help-title">Help</Typography>
                <IconButton onClick={onClose} aria-label="close help"><CloseIcon /></IconButton>
            </Box>
            <Divider />
            <Box sx={bodySx}>
                <Box component="nav" sx={railSx} aria-label="help pages">
                    <List>
                        {pages.map((p) => (
                            <ListItemButton key={p.key} selected={p.key === activeKey}
                                onClick={() => { setActiveKey(p.key); }}>
                                <ListItemText primary={p.title} />
                            </ListItemButton>
                        ))}
                    </List>
                </Box>
                <Box component="article" sx={contentSx}>
                    {Active ? <Active /> : null}
                </Box>
            </Box>
        </Drawer>
    );
};

export default HelpPanel;
