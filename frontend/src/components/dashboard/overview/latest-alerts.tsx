import * as React from 'react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardHeader from '@mui/material/CardHeader';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import ListItemText from '@mui/material/ListItemText';
import type { SxProps } from '@mui/material/styles';
import { ArrowRight as ArrowRightIcon } from '@phosphor-icons/react/dist/ssr/ArrowRight';
import { DotsThreeVertical as DotsThreeVerticalIcon } from '@phosphor-icons/react/dist/ssr/DotsThreeVertical';
import Avatar from '@mui/material/Avatar';
import dayjs from 'dayjs';

export interface Alert {
  id: string;
  icon?: React.ReactElement;
  name: string;
  updatedAt: Date;
}

export interface LatestAlertsProps {
  products?: Alert[];
  sx?: SxProps;
}

export function LatestAlerts({ products = [], sx }: LatestAlertsProps): React.JSX.Element {
  return (
    <Card
      sx={{
        display: 'flex',
        flexDirection: 'column',
        height: '100%',
        ...sx,
      }}
    >
      <CardHeader title="Últimas Alertas" />
      <Divider />
      
      <List sx={{ flexGrow: 1 }}>
        {products.length > 0 ? (
          products.map((product, index) => (
            <ListItem divider={index < products.length - 1} key={product.id}>
              <ListItemAvatar sx={{ mr: 2 }}>
                <Avatar
                  sx={{
                    backgroundColor: 'var(--mui-palette-neutral-100)',
                    color: 'var(--mui-palette-primary-main)',
                    height: 48,
                    width: 48,
                  }}
                >
                  {product.icon}
                </Avatar>
              </ListItemAvatar>
              <ListItemText
                primary={product.name}
                primaryTypographyProps={{ variant: 'subtitle1' }}
                secondary={`Actualizado ${dayjs(product.updatedAt).format('MMM D, YYYY - HH:MM')}`}
                secondaryTypographyProps={{ variant: 'body2' }}
              />
              <IconButton edge="end">
                <DotsThreeVerticalIcon weight="bold" />
              </IconButton>
            </ListItem>
          ))
        ) : (
          <ListItem>
            <ListItemText primary="No hay alertas" />
          </ListItem>
        )}
      </List>

      <Divider />
      <CardActions sx={{ justifyContent: 'flex-end', mt: 'auto' }}>
        <Button
          color="inherit"
          endIcon={<ArrowRightIcon fontSize="var(--icon-fontSize-md)" />}
          size="small"
          variant="text"
        >
          Ver todo
        </Button>
      </CardActions>
    </Card>

  );
}
