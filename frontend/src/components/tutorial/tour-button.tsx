'use client'; 

import { Box, IconButton, Tooltip } from "@mui/material";
import { HelpCircle as HelpIcon } from "lucide-react";
import React from "react";

interface TourButtonProps {
    onClick: () => void;
    tooltipTitle: string;
    style: string; //Este campo, colocar "none" si no se quiere animacion, si no se recomienda usar "pulse 3s infinite"
}

export const TourButton: React.FC<TourButtonProps> = ({ onClick, tooltipTitle, style }) => {
    return (
        <Tooltip title={tooltipTitle}>
            <Box sx={{
                animation: `${style}`, 
                borderRadius: '50%',
                display: 'flex',
            }}>
                <IconButton 
                    color="primary"
                    size="medium" 
                    onClick={onClick} 
                    aria-label="Iniciar tutorial"
                >
                    <HelpIcon size={24} /> 
                </IconButton>
            </Box>
        </Tooltip>
    );
};