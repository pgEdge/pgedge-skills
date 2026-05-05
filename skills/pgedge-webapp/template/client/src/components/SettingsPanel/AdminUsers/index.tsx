/*-------------------------------------------------------------------------
 *
 * <PROJECT_NAME>
 *
 * Copyright (c) <CURRENT_YEAR>, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

import { useState, useEffect, useCallback } from 'react';
import {
    Box, Typography, Button, Table, TableBody, TableCell, TableContainer,
    TableHead, TableRow, Chip, IconButton, Alert, CircularProgress, Paper,
} from '@mui/material';
import {
    Add as AddIcon,
    Edit as EditIcon,
    Key as KeyIcon,
    Delete as DeleteIcon,
} from '@mui/icons-material';
import { apiGet, apiDelete } from '../../../utils/apiClient';
import { useAuth } from '../../../contexts/useAuth';
import UserDialog from './UserDialog';
import AdminPasswordResetDialog from './AdminPasswordResetDialog';
import DeleteConfirmationDialog from './DeleteConfirmationDialog';

export interface UserListItem {
    username: string;
    full_name?: string;
    email?: string;
    is_superuser: boolean;
    enabled: boolean;
}

const AdminUsers: React.FC = () => {
    const { user } = useAuth();
    const [users, setUsers] = useState<UserListItem[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [addOpen, setAddOpen] = useState(false);
    const [editUser, setEditUser] = useState<UserListItem | null>(null);
    const [resetUser, setResetUser] = useState<UserListItem | null>(null);
    const [deleteUser, setDeleteUser] = useState<UserListItem | null>(null);

    const fetchUsers = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const data = await apiGet<UserListItem[]>('/api/v1/users');
            setUsers(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load users.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        if (user?.isSuperuser) {
            void fetchUsers();
        }
    }, [user?.isSuperuser, fetchUsers]);

    if (!user?.isSuperuser) {
        return (
            <Alert severity="warning">
                Superuser only — your account does not have administrative privileges.
            </Alert>
        );
    }

    const handleDeleteConfirm = async (username: string) => {
        await apiDelete(`/api/v1/users/${username}`);
        setDeleteUser(null);
        void fetchUsers();
    };

    return (
        <Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6" component="h2">Users</Typography>
                <Button
                    variant="contained"
                    startIcon={<AddIcon />}
                    onClick={() => { setAddOpen(true); }}
                >
                    Add User
                </Button>
            </Box>

            {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

            {loading ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
                    <CircularProgress aria-label="Loading users" />
                </Box>
            ) : (
                <TableContainer component={Paper} variant="outlined">
                    <Table aria-label="users table">
                        <TableHead>
                            <TableRow>
                                <TableCell>Username</TableCell>
                                <TableCell>Full Name</TableCell>
                                <TableCell>Email</TableCell>
                                <TableCell>Superuser</TableCell>
                                <TableCell>Enabled</TableCell>
                                <TableCell align="right">Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {users.map((u) => (
                                <TableRow key={u.username}>
                                    <TableCell>{u.username}</TableCell>
                                    <TableCell>{u.full_name ?? ''}</TableCell>
                                    <TableCell>{u.email ?? ''}</TableCell>
                                    <TableCell>
                                        <Chip
                                            label={u.is_superuser ? 'Yes' : 'No'}
                                            color={u.is_superuser ? 'primary' : 'default'}
                                            size="small"
                                        />
                                    </TableCell>
                                    <TableCell>
                                        <Chip
                                            label={u.enabled ? 'Enabled' : 'Disabled'}
                                            color={u.enabled ? 'success' : 'default'}
                                            size="small"
                                        />
                                    </TableCell>
                                    <TableCell align="right">
                                        <IconButton
                                            size="small"
                                            aria-label={`edit ${u.username}`}
                                            onClick={() => { setEditUser(u); }}
                                        >
                                            <EditIcon fontSize="small" />
                                        </IconButton>
                                        <IconButton
                                            size="small"
                                            aria-label={`reset password for ${u.username}`}
                                            onClick={() => { setResetUser(u); }}
                                        >
                                            <KeyIcon fontSize="small" />
                                        </IconButton>
                                        <IconButton
                                            size="small"
                                            aria-label={`delete ${u.username}`}
                                            onClick={() => { setDeleteUser(u); }}
                                        >
                                            <DeleteIcon fontSize="small" />
                                        </IconButton>
                                    </TableCell>
                                </TableRow>
                            ))}
                            {users.length === 0 && !loading && (
                                <TableRow>
                                    <TableCell colSpan={6} align="center">
                                        <Typography variant="body2" color="text.secondary">
                                            No users found.
                                        </Typography>
                                    </TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                </TableContainer>
            )}

            <UserDialog
                mode="add"
                open={addOpen}
                onClose={() => { setAddOpen(false); }}
                onSaved={() => { setAddOpen(false); void fetchUsers(); }}
            />

            {editUser && (
                <UserDialog
                    mode="edit"
                    initialUser={editUser}
                    open={Boolean(editUser)}
                    onClose={() => { setEditUser(null); }}
                    onSaved={() => { setEditUser(null); void fetchUsers(); }}
                />
            )}

            {resetUser && (
                <AdminPasswordResetDialog
                    username={resetUser.username}
                    open={Boolean(resetUser)}
                    onClose={() => { setResetUser(null); }}
                    onSaved={() => { setResetUser(null); }}
                />
            )}

            {deleteUser && (
                <DeleteConfirmationDialog
                    username={deleteUser.username}
                    open={Boolean(deleteUser)}
                    onClose={() => { setDeleteUser(null); }}
                    onConfirm={handleDeleteConfirm}
                />
            )}
        </Box>
    );
};

export default AdminUsers;
