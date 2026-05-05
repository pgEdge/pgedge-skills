/*-------------------------------------------------------------------------
 *
 * <PROJECT_NAME>
 *
 * Copyright (c) <CURRENT_YEAR>, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

import { useState } from 'react';
import {
    Dialog, DialogTitle, DialogContent, DialogActions,
    Button, TextField, Alert, Box,
} from '@mui/material';
import type { Theme } from '@mui/material/styles';
import { SELECT_FIELD_SX } from '../shared/formStyles';
import { apiPost } from '../../utils/apiClient';

interface ChangePasswordDialogProps {
    open: boolean;
    onClose: () => void;
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

const ChangePasswordDialog: React.FC<ChangePasswordDialogProps> = ({ open, onClose }) => {
    const [currentPassword, setCurrentPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [submitting, setSubmitting] = useState(false);

    const handleClose = () => {
        setCurrentPassword('');
        setNewPassword('');
        setConfirmPassword('');
        setError(null);
        setSubmitting(false);
        onClose();
    };

    const handleSubmit = async () => {
        setError(null);

        if (!currentPassword) {
            setError('Current password is required.');
            return;
        }
        if (!newPassword) {
            setError('New password is required.');
            return;
        }
        if (newPassword !== confirmPassword) {
            setError('New password and confirmation do not match.');
            return;
        }

        setSubmitting(true);
        try {
            await apiPost('/api/v1/user/password', {
                current_password: currentPassword,
                new_password: newPassword,
            });
            handleClose();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to change password.');
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <Dialog
            open={open}
            onClose={handleClose}
            aria-labelledby="change-password-dialog-title"
            maxWidth="sm"
            fullWidth
        >
            <DialogTitle id="change-password-dialog-title">Change Password</DialogTitle>
            <DialogContent>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
                    {error && (
                        <Alert severity="error">{error}</Alert>
                    )}
                    <TextField
                        label="Current password"
                        type="password"
                        value={currentPassword}
                        onChange={(e) => { setCurrentPassword(e.target.value); }}
                        InputLabelProps={{ shrink: true }}
                        sx={(theme) => ({ ...textFieldSx(theme), ...SELECT_FIELD_SX })}
                        fullWidth
                        autoComplete="current-password"
                    />
                    <TextField
                        label="New password"
                        type="password"
                        value={newPassword}
                        onChange={(e) => { setNewPassword(e.target.value); }}
                        InputLabelProps={{ shrink: true }}
                        sx={(theme) => ({ ...textFieldSx(theme), ...SELECT_FIELD_SX })}
                        fullWidth
                        autoComplete="new-password"
                    />
                    <TextField
                        label="Confirm password"
                        type="password"
                        value={confirmPassword}
                        onChange={(e) => { setConfirmPassword(e.target.value); }}
                        InputLabelProps={{ shrink: true }}
                        sx={(theme) => ({ ...textFieldSx(theme), ...SELECT_FIELD_SX })}
                        fullWidth
                        autoComplete="new-password"
                    />
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose} disabled={submitting}>Cancel</Button>
                <Button
                    onClick={() => { void handleSubmit(); }}
                    variant="contained"
                    disabled={submitting}
                >
                    {submitting ? 'Saving…' : 'Change Password'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default ChangePasswordDialog;
