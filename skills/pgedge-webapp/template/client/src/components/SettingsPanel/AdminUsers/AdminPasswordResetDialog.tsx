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
    Button, TextField, Alert, Box,
} from '@mui/material';
import type { Theme } from '@mui/material/styles';
import { SELECT_FIELD_SX } from '../../shared/formStyles';
import { apiPost } from '../../../utils/apiClient';

interface AdminPasswordResetDialogProps {
    username: string;
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

const AdminPasswordResetDialog: React.FC<AdminPasswordResetDialogProps> = ({
    username, open, onClose, onSaved,
}) => {
    const [newPassword, setNewPassword] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        if (open) {
            setNewPassword('');
            setError(null);
            setSubmitting(false);
        }
    }, [open]);

    const handleClose = () => {
        onClose();
    };

    const handleSubmit = async () => {
        setError(null);

        if (!newPassword) {
            setError('New password is required.');
            return;
        }

        setSubmitting(true);
        try {
            await apiPost(`/api/v1/users/${username}/password`, {
                new_password: newPassword,
            });
            onSaved();
            onClose();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to reset password.');
        } finally {
            setSubmitting(false);
        }
    };

    const titleId = 'admin-password-reset-dialog-title';

    return (
        <Dialog
            open={open}
            onClose={handleClose}
            aria-labelledby={titleId}
            maxWidth="sm"
            fullWidth
        >
            <DialogTitle id={titleId}>Reset Password: {username}</DialogTitle>
            <DialogContent>
                <Box sx={{ pt: 1 }}>
                    {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
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
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose} disabled={submitting}>Cancel</Button>
                <Button
                    onClick={() => { void handleSubmit(); }}
                    variant="contained"
                    disabled={submitting}
                >
                    {submitting ? 'Saving…' : 'Reset Password'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default AdminPasswordResetDialog;
