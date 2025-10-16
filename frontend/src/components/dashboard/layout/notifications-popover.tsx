/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
'use client';


import * as React from 'react';
import Box from '@mui/material/Box';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import ListItemText from '@mui/material/ListItemText';
import Avatar from '@mui/material/Avatar';
import Popover from '@mui/material/Popover';
import Typography from '@mui/material/Typography';
import Chip from '@mui/material/Chip';
import Tooltip from '@mui/material/Tooltip';
import { Trash as TrashIcon } from '@phosphor-icons/react/dist/ssr/Trash';
import { Check as CheckIcon } from '@phosphor-icons/react/dist/ssr/Check';
import { Alarm as AlarmIcon } from '@phosphor-icons/react/dist/ssr/Alarm';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { useNotifications } from '@/contexts/notifications';


dayjs.extend(relativeTime);


export interface NotificationsPopoverProps {
    anchorEl: Element | null;
    onClose: () => void;
    open: boolean;
}


export function NotificationsPopover({ anchorEl, onClose, open }: NotificationsPopoverProps) {
    const { notifications, markAllRead, remove, clear } = useNotifications();


    return (
        <Popover
            anchorEl={anchorEl}
            anchorOrigin={{ horizontal: 'left', vertical: 'bottom' }}
            onClose={onClose}
            open={open}
            slotProps={{ paper: { sx: { width: 380, maxHeight: 500 } } }}
        >
        <Box sx={{ p: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
            <Typography variant="subtitle1">Notificaciones</Typography>
            <Chip size="small" label={`${notifications.length}`} />
            <Box sx={{ flex: 1 }} />
            <Tooltip title="Marcar todas como leídas">
                <IconButton onClick={markAllRead}>
                    <CheckIcon />
                </IconButton>
            </Tooltip>
            <Tooltip title="Eliminar todas">
                <IconButton onClick={clear}>
                    <TrashIcon />
                </IconButton>
            </Tooltip>
        </Box>
        <Divider />
        <List dense disablePadding>
            {notifications.length === 0 && (
                <Box sx={{ p: 2 }}>
                    <Typography color="text.secondary" variant="body2">
                    Sin notificaciones por ahora.
                    </Typography>
                </Box>
            )}
            {notifications.map((n) => (
                <ListItem key={n.id} secondaryAction={
                    <Tooltip title="Eliminar">
                        <IconButton edge="end" onClick={() => { remove(n.id); }}>
                            <TrashIcon />
                        </IconButton>
                    </Tooltip>
                }>
                <ListItemAvatar>
                    <Avatar>
                        <AlarmIcon />
                    </Avatar>
                </ListItemAvatar>
                <ListItemText
                    primary={n.title}
                    secondary={
                    <>
                        <Typography component="span" variant="body2" color="text.secondary">
                            {n.body}
                        </Typography>
                        <br />
                        <Typography component="span" variant="caption" color="text.secondary">
                            {dayjs(n.ts).fromNow()}
                        </Typography>
                    </>
                    }
                />
            </ListItem>
            ))}
        </List>
        </Popover>
    );
}
