import type { Icon } from '@phosphor-icons/react/dist/lib/types';
import { ChartPie as ChartPieIcon } from '@phosphor-icons/react/dist/ssr/ChartPie';
import { GearSix as GearSixIcon } from '@phosphor-icons/react/dist/ssr/GearSix';
import { PlugsConnected as PlugsConnectedIcon } from '@phosphor-icons/react/dist/ssr/PlugsConnected';
import { User as UserIcon } from '@phosphor-icons/react/dist/ssr/User';
import { Users as UsersIcon } from '@phosphor-icons/react/dist/ssr/Users';
import { XSquare } from '@phosphor-icons/react/dist/ssr/XSquare';
import { WarningCircleIcon } from '@phosphor-icons/react/dist/ssr/WarningCircle';
import { BuildingApartmentIcon } from '@phosphor-icons/react/dist/ssr/BuildingApartment';
import { ChartLineIcon } from '@phosphor-icons/react/dist/ssr/ChartLine';
import { FileTextIcon } from '@phosphor-icons/react/dist/ssr/FileText';
import { BlueprintIcon } from '@phosphor-icons/react/dist/ssr/Blueprint';

export const navIcons = {
  'chart-pie': ChartPieIcon,
  'gear-six': GearSixIcon,
  'plugs-connected': PlugsConnectedIcon,
  'x-square': XSquare,
  'warning-circle': WarningCircleIcon,
  'building-apartment': BuildingApartmentIcon,
  'chart-line': ChartLineIcon,
  'file-text': FileTextIcon,
  'blueprint': BlueprintIcon,
  user: UserIcon,
  users: UsersIcon,
} as Record<string, Icon>;
