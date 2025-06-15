import type { Icon } from '@phosphor-icons/react/dist/lib/types';
import { ChartPieIcon as ChartPieIcon } from '@phosphor-icons/react/dist/ssr/ChartPie';
import { GearSixIcon as GearSixIcon } from '@phosphor-icons/react/dist/ssr/GearSix';
import { PlugsConnectedIcon as PlugsConnectedIcon } from '@phosphor-icons/react/dist/ssr/PlugsConnected';
import { UserIcon as UserIcon } from '@phosphor-icons/react/dist/ssr/User';
import { UsersIcon as UsersIcon } from '@phosphor-icons/react/dist/ssr/Users';
import { XSquareIcon as XSquare } from '@phosphor-icons/react/dist/ssr/XSquare';
import { WarningOctagonIcon as Warning } from '@phosphor-icons/react/dist/ssr/WarningOctagon';

//Iconos de: https://phosphoricons.com/

export const navIcons = {
  'chart-pie': ChartPieIcon,
  'gear-six': GearSixIcon,
  'plugs-connected': PlugsConnectedIcon,
  'x-square': XSquare,
  'warning': Warning,
  user: UserIcon,
  users: UsersIcon,
} as Record<string, Icon>;
