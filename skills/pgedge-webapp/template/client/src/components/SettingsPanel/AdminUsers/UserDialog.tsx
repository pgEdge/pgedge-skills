/*-------------------------------------------------------------------------
 *
 * <PROJECT_NAME>
 *
 * Copyright (c) <CURRENT_YEAR>, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

import { useState, useEffect } from 'react';
import {
    Dialog, DialogTitle, DialogContent, DialogActions,
    Button, TextField, Alert, Box, Switch, FormControlLabel,
} from '@mui/material';
import type { Theme } from '@mui/material/styles';
import { SELECT_FIELD_SX } from '../../shared/formStyles';
import { apiPost, apiPatch } from '../../../utils/apiClient';
import type { UserListItem } from './index';

interface UserDialogProps {
    mode: 'add' | 'edit';
    initialUser?: UserListItem;
    open: boolean;
    onClose: () => void;
    onSaved: () => void;
}

const textFieldSx = (theme: Theme) => ({
    '& .MuiOutlinedInput-root': {
        '&:hover .MuiOutlinedInput-notchedOutline': { borderColor: theme.palette.grey[400] },
        '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
            borderColor: theme.palette.primary.main, borderWidth: 2,
        },
    },
    '& .MuiInputLabel-root.Mui-focused': { color: 'primary.main' },
});

const UserDialog: React.FC<UserDialogProps> = ({
    mode, initialUser, open, onClose, onSaved,
}) => {
    const [username, setUsername] = useState('');
    const [fullName, setFullName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [isSuperuser, setIsSuperuser] = useState(false);
    const [enabled, setEnabled] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        if (open) {
            if (mode === 'edit' && initialUser) {
                setUsername(initialUser.username);
                setFullName(initialUser.full_name ?? '');
                setEmail(initialUser.email ?? '');
                setIsSuperuser(initialUser.is_superuser);
                setEnabled(initialUser.enabled);
            } else {
                setUsername('');
                setFullName('');
                setEmail('');
                setPassword('');
                setIsSuperuser(false);
                setEnabled(true);
            }
            setError(null);
            setSubmitting(false);
        }
    }, [open, mode, initialUser]);

    const handleClose = () => {
        onClose();
    };

    const handleSubmit = async () => {
        setError(null);

        if (mode === 'add') {
            if (!username.trim()) {
                setError('Username is required.');
                return;
            }
            if (!password) {
                setError('Password is required.');
                return;
            }
        }

        setSubmitting(true);
        try {
            if (mode === 'add') {
                await apiPost('/api/v1/users', {
                    username: username.trim(),
                    password,
                    full_name: fullName || undefined,
                    email: email || undefined,
                    is_superuser: isSuperuser,
                });
            } else {
                await apiPatch(`/api/v1/users/${initialUser!.username}`, {
                    full_name: fullName || undefined,
                    email: email || undefined,
                    is_superuser: isSuperuser,
                    enabled,
                });
            }
            onSaved();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to save user.');
        } finally {
            setSubmitting(false);
        }
    };

    const titleId = 'user-dialog-title';
    const title = mode === 'add' ? 'Add User' : `Edit User: ${initialUser?.username ?? ''}`;

    return (
        <Dialog
            open={open}
            onClose={handleClose}
            aria-labelledby={titleId}
            maxWidth="sm"
            fullWidth
        >
            <DialogTitle id={titleId}>{title}</DialogTitle>
            <DialogContent>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
                    {error && <Alert severity="error">{error}</Alert>}

                    <TextField
                        label="Username"
                        value={username}
                        onChange={(e) => { setUsername(e.target.value); }}
                        InputLabelProps={{ shrink: true }}
                        sx={(theme) => ({ ...textFieldSx(theme), ...SELECT_FIELD_SX })}
                        fullWidth
                        disabled={mode === 'edit'}
                        inputProps={mode === 'edit' ? { readOnly: true } : undefined}
                    />

                    <TextField
                        label="Full Name"
                        value={fullName}
                        onChange={(e) => { setFullName(e.target.value); }}
                        InputLabelProps={{ shrink: true }}
                        sx={(theme) => ({ ...textFieldSx(theme), ...SELECT_FIELD_SX })}
                        fullWidth
                    />

                    <TextField
                        label="Email"
                        type="email"
                        value={email}
                        onChange={(e) => { setEmail(e.target.value); }}
                        InputLabelProps={{ shrink: true }}
                        sx={(theme) => ({ ...textFieldSx(theme), ...SELECT_FIELD_SX })}
                        fullWidth
                    />

                    {mode === 'add' && (
                        <TextField
                            label="Password"
                            type="password"
                            value={password}
                            onChange={(e) => { setPassword(e.target.value); }}
                            InputLabelProps={{ shrink: true }}
                            sx={(theme) => ({ ...textFieldSx(theme), ...SELECT_FIELD_SX })}
                            fullWidth
                            autoComplete="new-password"
                        />
                    )}

                    <FormControlLabel
                        control={
                            <Switch
                                checked={isSuperuser}
                                onChange={(e) => { setIsSuperuser(e.target.checked); }}
                            />
                        }
                        label="Is Superuser"
                    />

                    {mode === 'edit' && (
                        <FormControlLabel
                            control={
                                <Switch
                                    checked={enabled}
                                    onChange={(e) => { setEnabled(e.target.checked); }}
                                />
                            }
                            label="Enabled"
                        />
                    )}
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose} disabled={submitting}>Cancel</Button>
                <Button
                    onClick={() => { void handleSubmit(); }}
                    variant="contained"
                    disabled={submitting}
                >
                    {submitting ? 'Saving…' : mode === 'add' ? 'Add User' : 'Save Changes'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default UserDialog;
