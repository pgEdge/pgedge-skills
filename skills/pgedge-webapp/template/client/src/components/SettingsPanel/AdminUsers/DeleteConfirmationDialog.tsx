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
    Dialog, DialogTitle, DialogContent, DialogContentText,
    DialogActions, Button, Alert,
} from '@mui/material';

interface DeleteConfirmationDialogProps {
    username: string;
    open: boolean;
    onClose: () => void;
    onConfirm: (username: string) => Promise<void>;
}

const DeleteConfirmationDialog: React.FC<DeleteConfirmationDialogProps> = ({
    username, open, onClose, onConfirm,
}) => {
    const [error, setError] = useState<string | null>(null);
    const [submitting, setSubmitting] = useState(false);

    const handleConfirm = async () => {
        setError(null);
        setSubmitting(true);
        try {
            await onConfirm(username);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to delete user.');
            setSubmitting(false);
        }
    };

    const handleClose = () => {
        setError(null);
        setSubmitting(false);
        onClose();
    };

    const titleId = 'delete-confirmation-dialog-title';

    return (
        <Dialog
            open={open}
            onClose={handleClose}
            aria-labelledby={titleId}
            maxWidth="xs"
            fullWidth
        >
            <DialogTitle id={titleId}>Delete User</DialogTitle>
            <DialogContent>
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
                <DialogContentText>
                    Are you sure you want to delete user <strong>{username}</strong>?
                    This action cannot be undone.
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose} disabled={submitting}>Cancel</Button>
                <Button
                    onClick={() => { void handleConfirm(); }}
                    variant="contained"
                    color="error"
                    disabled={submitting}
                >
                    {submitting ? 'Deleting…' : 'Delete'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default DeleteConfirmationDialog;
