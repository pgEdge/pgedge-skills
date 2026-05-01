import React, { useMemo, useState } from 'react';
import {
    Dialog, AppBar, Toolbar, IconButton, Typography, Box,
    List, ListItemButton, ListItemText,
} from '@mui/material';
import { useTheme } from '@mui/material/styles';
import { Close as CloseIcon } from '@mui/icons-material';
import SlideTransition from '../shared/SlideTransition';

interface SettingsItem {
    id: string;
    label: string;
    Component: React.FC;
}

interface SettingsSection {
    category: string;
    items: SettingsItem[];
}

const GeneralPlaceholder: React.FC = () => (
    <Typography variant="body1">
        This is where settings go. Add sections and items as your application
        grows by extending the SECTIONS array in this file.
    </Typography>
);

/**
 * Replace or extend the entries below to add settings sections. Each item
 * renders its Component in the right-hand pane when selected.
 */
const SECTIONS: SettingsSection[] = [
    {
        category: 'General',
        items: [
            { id: 'general', label: 'General', Component: GeneralPlaceholder },
        ],
    },
];

interface SettingsPanelProps {
    open: boolean;
    onClose: () => void;
}

const subsectionLabelSx = {
    textTransform: 'uppercase' as const,
    fontSize: '0.7rem',
    fontWeight: 700,
    letterSpacing: '0.06em',
    color: 'text.secondary',
};

const SettingsPanel: React.FC<SettingsPanelProps> = ({ open, onClose }) => {
    const theme = useTheme();
    const [activeId, setActiveId] = useState<string>('');

    const allItems = useMemo(() => SECTIONS.flatMap((s) => s.items), []);

    const handleEnter = () => {
        if (allItems.length > 0) {
            setActiveId(allItems[0].id);
        }
    };

    const activeItem = allItems.find((item) => item.id === activeId);
    const ActiveComponent = activeItem?.Component;

    return (
        <Dialog
            fullScreen
            open={open}
            onClose={onClose}
            TransitionComponent={SlideTransition}
            TransitionProps={{ onEnter: handleEnter }}
            aria-labelledby="settings-title"
        >
            <AppBar
                position="static"
                elevation={0}
                sx={{
                    bgcolor: theme.palette.background.paper,
                    borderBottom: '1px solid',
                    borderColor: theme.palette.divider,
                }}
            >
                <Toolbar>
                    <IconButton
                        edge="start"
                        color="inherit"
                        onClick={onClose}
                        aria-label="close settings"
                        sx={{ color: theme.palette.text.secondary, mr: 2 }}
                    >
                        <CloseIcon />
                    </IconButton>
                    <Typography
                        id="settings-title"
                        variant="h6"
                        component="div"
                        sx={{
                            flexGrow: 1,
                            fontWeight: 600,
                            color: theme.palette.text.primary,
                        }}
                    >
                        Settings
                    </Typography>
                </Toolbar>
            </AppBar>
            <Box sx={{ display: 'flex', flex: 1, overflow: 'hidden' }}>
                <Box
                    component="nav"
                    aria-label="settings sections"
                    sx={{
                        width: 240,
                        flexShrink: 0,
                        borderRight: '1px solid',
                        borderColor: theme.palette.divider,
                        bgcolor: theme.palette.background.paper,
                        overflowY: 'auto',
                    }}
                >
                    <List disablePadding sx={{ py: 1 }}>
                        {SECTIONS.map((section) => (
                            <React.Fragment key={section.category}>
                                <Typography sx={{ ...subsectionLabelSx, px: 2, pt: 2, pb: 0.5 }}>
                                    {section.category}
                                </Typography>
                                {section.items.map((item) => {
                                    const isSelected = item.id === activeId;
                                    return (
                                        <ListItemButton
                                            key={item.id}
                                            selected={isSelected}
                                            onClick={() => { setActiveId(item.id); }}
                                            sx={{
                                                borderRadius: 1,
                                                mx: 1,
                                                bgcolor: isSelected
                                                    ? theme.palette.action.selected
                                                    : 'transparent',
                                            }}
                                        >
                                            <ListItemText
                                                primary={item.label}
                                                primaryTypographyProps={{ fontSize: '1rem' }}
                                            />
                                        </ListItemButton>
                                    );
                                })}
                            </React.Fragment>
                        ))}
                    </List>
                </Box>
                <Box
                    component="main"
                    sx={{
                        flex: 1,
                        overflow: 'auto',
                        bgcolor: theme.palette.background.default,
                        p: 3,
                    }}
                >
                    {ActiveComponent && <ActiveComponent />}
                </Box>
            </Box>
        </Dialog>
    );
};

export default SettingsPanel;
