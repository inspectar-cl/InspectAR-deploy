import type { Icon } from '@phosphor-icons/react/dist/lib/types';
import { ChartPieIcon } from '@phosphor-icons/react/dist/ssr/ChartPie';
import { BroadcastIcon } from '@phosphor-icons/react/dist/ssr/Broadcast';
import { GearSixIcon } from '@phosphor-icons/react/dist/ssr/GearSix';
import { PlugsConnectedIcon } from '@phosphor-icons/react/dist/ssr/PlugsConnected';
import { UserIcon } from '@phosphor-icons/react/dist/ssr/User';
import { UsersIcon } from '@phosphor-icons/react/dist/ssr/Users';
import { XSquareIcon as XSquare } from '@phosphor-icons/react/dist/ssr/XSquare';
import { WarningCircleIcon } from '@phosphor-icons/react/dist/ssr/WarningCircle';
import { BuildingApartmentIcon } from '@phosphor-icons/react/dist/ssr/BuildingApartment';
import { ChartLineIcon } from '@phosphor-icons/react/dist/ssr/ChartLine';
import { FileTextIcon } from '@phosphor-icons/react/dist/ssr/FileText';
import { BlueprintIcon } from '@phosphor-icons/react/dist/ssr/Blueprint';
import { AddressBookIcon } from '@phosphor-icons/react/dist/ssr/AddressBook';
import { SirenIcon } from '@phosphor-icons/react/dist/ssr/Siren';
import { BrainIcon } from '@phosphor-icons/react/dist/ssr/Brain';
import { ClipboardTextIcon } from '@phosphor-icons/react/dist/ssr/ClipboardText';
import { CloudArrowUpIcon } from '@phosphor-icons/react/dist/ssr/CloudArrowUp';

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
  'broadcast': BroadcastIcon,
  'adressBook': AddressBookIcon,
  'siren': SirenIcon,
  'brain': BrainIcon,
  'clipboard-text': ClipboardTextIcon,
  'cloud-arrow-up': CloudArrowUpIcon,
  user: UserIcon,
  users: UsersIcon,
} as Record<string, Icon>;
