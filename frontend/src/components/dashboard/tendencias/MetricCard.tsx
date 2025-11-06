import { Box, Card, CardContent, Typography } from '@mui/material';
import { TrendingUp, TrendingDown } from 'lucide-react';
import type { ReactNode } from 'react';

interface MetricCardProps {
  title: string;
  value: number | string;
  color: string;
  icon: ReactNode;
  suffix?: string;
  trend?: {
    value: number;
    isPositive: boolean;
  };
}

export default function MetricCard({ title, value, color, icon, suffix = '', trend }: MetricCardProps) {
  return (
    <Card
      sx={{
        height: '100%',
        borderLeft: `4px solid ${color}`,
        transition: 'transform 0.2s, box-shadow 0.2s',
        '&:hover': {
          transform: 'translateY(-4px)',
          boxShadow: 'var(--mui-shadows-8)'
        }
      }}
    >
      <CardContent>
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={1}>
          <Typography variant="body2" color="text.secondary" fontWeight={500}>
            {title}
          </Typography>
          <Box sx={{ color, opacity: 0.8 }}>
            {icon}
          </Box>
        </Box>

        <Typography variant="h4" fontWeight="bold" mb={0.5}>
          {value}{suffix}
        </Typography>

        {trend && (
          <Box display="flex" alignItems="center" gap={0.5}>
            {trend.isPositive ? (
              <TrendingUp size={16} color="var(--mui-palette-success-main)" />
            ) : (
              <TrendingDown size={16} color="var(--mui-palette-error-main)" />
            )}
            <Typography
              variant="caption"
              sx={{
                color: trend.isPositive ? 'var(--mui-palette-success-main)' : 'var(--mui-palette-error-main)',
                fontWeight: 500
              }}
            >
              {trend.value}% vs mes anterior
            </Typography>
          </Box>
        )}
      </CardContent>
    </Card>
  );
}