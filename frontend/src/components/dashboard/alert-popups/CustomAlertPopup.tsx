'use client';

import * as React from 'react';
import Box from '@mui/material/Box';
import Alert from '@mui/material/Alert';
import AlertTitle from '@mui/material/AlertTitle';
import IconButton from '@mui/material/IconButton';
import Collapse from '@mui/material/Collapse';
import CloseIcon from '@mui/icons-material/Close';

interface CustomAlertProps {
  title: string;
  message: string;
  severity: 'success' | 'info' | 'warning' | 'error';
  open: boolean;
  onClose: () => void;
}

export function CustomAlert({ title, message, severity, open, onClose }: CustomAlertProps): React.JSX.Element {
  return (
    <Box 
      sx={{ 
        position: 'fixed', 
        top: 20, 
        left: '55%', 
        transform: 'translateX(-50%)', 
        zIndex: 9999,
        width: 'auto',
        maxWidth: 450,
      }}
    >
      <Collapse in={open}>
        <Alert
          severity={severity}
          action={
            <IconButton
              aria-label="close"
              color="inherit"
              size="small"
              onClick={onClose}
            >
              <CloseIcon fontSize="inherit" />
            </IconButton>
          }
        >
          <AlertTitle>{title}</AlertTitle>
          {message}
        </Alert>
      </Collapse>
    </Box>
  );
}