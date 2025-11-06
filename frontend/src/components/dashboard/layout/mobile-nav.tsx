/* eslint-disable @typescript-eslint/no-unsafe-argument -- Navigation config requires dynamic typing */
'use client';

import * as React from 'react';
import RouterLink from 'next/link';
import { usePathname } from 'next/navigation';
import Box from '@mui/material/Box';
import Divider from '@mui/material/Divider';
import Drawer from '@mui/material/Drawer';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import ButtonBase from '@mui/material/ButtonBase';
import { CaretUpDownIcon } from '@phosphor-icons/react/dist/ssr/CaretUpDown';

import type { NavItemConfig } from '@/types/nav';
import { paths } from '@/paths';
import { isNavItemActive } from '@/lib/is-nav-item-active';
import { Logo } from '@/components/core/logo';

import { navItems } from './config';
import { navIcons } from './nav-icons';

import { rolePermissions } from '@/role-permissions';
import { useAuthUser } from '@/contexts/user-context';

interface Edificio {
    id: number;
    nombre: string;
    direccion: string;
    creado_en: string;
}

export interface MobileNavProps {
  onClose?: () => void;
  open?: boolean;
  items?: NavItemConfig[];
}

export function MobileNav({ open, onClose }: MobileNavProps): React.JSX.Element {
  const pathname = usePathname();
  
    const { user, changeSelectedEdificio } = useAuthUser();
    const role = user?.role || 'residente';
    const allowedPaths = rolePermissions[role] || [];
  
    const filteredItems = navItems.filter((item) =>
      item.href ? allowedPaths.includes(item.href) : true
    );

    const selectedEdificio = user?.selected_edificio;
    const edificioNombre = selectedEdificio?.nombre ?? 'Cargando Edificio...';
    const edificiosList = user?.edificio;

    const isSelectable = edificiosList && edificiosList.length > 1;

    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const menuOpen = Boolean(anchorEl);

    const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        if (edificiosList && edificiosList.length > 1) { 
            setAnchorEl(event.currentTarget); 
        }
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const handleSelect = (edificio: Edificio) => {
        if (user?.selected_edificio?.id !== edificio.id) {
            void changeSelectedEdificio(edificio); 
        }
        handleClose();
        onClose?.(); //Esto cierra el menu cuando se selecciona edificio
    }

  return (
    <Drawer
      PaperProps={{
        sx: {
          '--MobileNav-background': 'var(--mui-palette-neutral-950)',
          '--MobileNav-color': 'var(--mui-palette-common-white)',
          '--NavItem-color': 'var(--mui-palette-neutral-300)',
          '--NavItem-hover-background': 'rgba(255, 255, 255, 0.04)',
          '--NavItem-active-background': 'var(--mui-palette-primary-main)',
          '--NavItem-active-color': 'var(--mui-palette-primary-contrastText)',
          '--NavItem-disabled-color': 'var(--mui-palette-neutral-500)',
          '--NavItem-icon-color': 'var(--mui-palette-neutral-400)',
          '--NavItem-icon-active-color': 'var(--mui-palette-primary-contrastText)',
          '--NavItem-icon-disabled-color': 'var(--mui-palette-neutral-600)',
          bgcolor: 'var(--MobileNav-background)',
          color: 'var(--MobileNav-color)',
          display: 'flex',
          flexDirection: 'column',
          maxWidth: '100%',
          scrollbarWidth: 'none',
          width: 'var(--MobileNav-width)',
          zIndex: 'var(--MobileNav-zIndex)',
          '&::-webkit-scrollbar': { display: 'none' },
        },
      }}
      onClose={onClose}
      open={open}
    >
      <Stack spacing={2} sx={{ p: 3 }}>
        <Box component={RouterLink} href={paths.home} sx={{ display: 'inline-flex' }}>
          <Logo color="light" height={32} width={122} />
        </Box>
        <Box 
          component={ButtonBase}
          onClick={handleClick}
          sx={{
              alignItems: 'center',
              backgroundColor: 'var(--mui-palette-neutral-950)',
              border: '1px solid var(--mui-palette-neutral-700)',
              borderRadius: '12px',
              cursor: isSelectable ? 'pointer' : 'default',
              display: 'flex',
              p: '4px 12px',
              textAlign: 'left',
              '&:hover': {
                  backgroundColor: isSelectable ? 'rgba(255, 255, 255, 0.04)' : undefined,
              }
          }}
          aria-controls={menuOpen ? 'edificio-menu-mobile' : undefined}
          aria-haspopup="true"
          aria-expanded={menuOpen ? 'true' : undefined}
        >
          <Box sx={{ flex: '1 1 auto', overflow: 'hidden' }}>
              <Typography color="var(--mui-palette-neutral-400)" variant="body2">
                  Edificio
              </Typography>
              <Typography 
                  color="inherit" 
                  variant="subtitle1"
                  sx={{
                      whiteSpace: 'nowrap',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                  }}
              >
                  {edificioNombre}
              </Typography>
          </Box>
          {isSelectable ? <CaretUpDownIcon /> : null} 
        </Box>

        <Menu
            anchorEl={anchorEl}
            id="edificio-menu-mobile"
            open={menuOpen}
            onClose={handleClose}
            slotProps={{
                paper: {
                    style: { maxHeight: 240 },
                    sx: {
                        bgcolor: 'var(--mui-palette-neutral-950)',
                        border: '1px solid var(--mui-palette-neutral-700)',
                        color: 'var(--MobileNav-color)',
                        borderRadius: '8px',
                        boxShadow: '0px 4px 12px rgba(0, 0, 0, 0.4)'
                    },
                },
                list: {
                    sx: { overflowY: 'auto', padding: 0 },
                },
            }}
            anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
            transformOrigin={{ vertical: 'top', horizontal: 'left' }}
          >
          {edificiosList?.map((edificio) => (
              <MenuItem 
                  key={edificio.id} 
                  onClick={() => {handleSelect(edificio)}}
                  selected={selectedEdificio?.id === edificio.id}
                  sx={{
                      borderRadius: '6px',
                      margin: '4px 4px',
                      padding: '8px 16px',
                      '&:first-of-type': { marginTop: '8px' },
                      '&:last-child': { marginBottom: '8px' },
                      '&:hover': { backgroundColor: 'rgba(255, 255, 255, 0.04)' },
                      '&.Mui-selected': {
                          backgroundColor: 'var(--mui-palette-primary-main)',
                          color: 'var(--mui-palette-primary-contrastText)',
                          '&:hover': {
                              backgroundColor: 'var(--mui-palette-primary-main)',
                          }
                      },
                      overflow: 'hidden',
                      '& > div': {
                          whiteSpace: 'nowrap',
                          overflow: 'hidden',
                          textOverflow: 'ellipsis',
                      }
                  }}
              >
                  {edificio.nombre}
              </MenuItem>
          ))}
          {edificiosList?.length === 0 && (
              <MenuItem disabled>No hay edificios disponibles</MenuItem>
          )}
      </Menu>
      </Stack>
      <Divider sx={{ borderColor: 'var(--mui-palette-neutral-700)' }} />
      <Box component="nav" sx={{ flex: '1 1 auto', p: '12px' }}>
        {renderNavItems({ pathname, items: filteredItems })}
      </Box>
    </Drawer>
  );
}

function renderNavItems({ items = [], pathname }: { items?: NavItemConfig[]; pathname: string }): React.JSX.Element {
  const children = items.reduce((acc: React.ReactNode[], curr: NavItemConfig): React.ReactNode[] => {
    const { key, ...item } = curr;

    acc.push(<NavItem key={key} pathname={pathname} {...item} />);

    return acc;
  }, []);

  return (
    <Stack component="ul" spacing={1} sx={{ listStyle: 'none', m: 0, p: 0 }}>
      {children}
    </Stack>
  );
}

interface NavItemProps extends Omit<NavItemConfig, 'items'> {
  pathname: string;
}

function NavItem({ disabled, external, href, icon, matcher, pathname, title }: NavItemProps): React.JSX.Element {
  const active = isNavItemActive({ disabled, external, href, matcher, pathname });
  const Icon = icon ? navIcons[icon] : null;

  return (
    <li>
      <Box
        {...(href
          ? {
              component: external ? 'a' : RouterLink,
              href,
              target: external ? '_blank' : undefined,
              rel: external ? 'noreferrer' : undefined,
            }
          : { role: 'button' })}
        sx={{
          alignItems: 'center',
          borderRadius: 1,
          color: 'var(--NavItem-color)',
          cursor: 'pointer',
          display: 'flex',
          flex: '0 0 auto',
          gap: 1,
          p: '6px 16px',
          position: 'relative',
          textDecoration: 'none',
          whiteSpace: 'nowrap',
          ...(disabled && {
            bgcolor: 'var(--NavItem-disabled-background)',
            color: 'var(--NavItem-disabled-color)',
            cursor: 'not-allowed',
          }),
          ...(active && { bgcolor: 'var(--NavItem-active-background)', color: 'var(--NavItem-active-color)' }),
        }}
      >
        <Box sx={{ alignItems: 'center', display: 'flex', justifyContent: 'center', flex: '0 0 auto' }}>
          {Icon ? (
            <Icon
              fill={active ? 'var(--NavItem-icon-active-color)' : 'var(--NavItem-icon-color)'}
              fontSize="var(--icon-fontSize-md)"
              weight={active ? 'fill' : undefined}
            />
          ) : null}
        </Box>
        <Box sx={{ flex: '1 1 auto' }}>
          <Typography
            component="span"
            sx={{ color: 'inherit', fontSize: '0.875rem', fontWeight: 500, lineHeight: '28px' }}
          >
            {title}
          </Typography>
        </Box>
      </Box>
    </li>
  );
}
