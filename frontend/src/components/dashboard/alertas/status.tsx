import * as React from 'react';
import Chip from '@mui/material/Chip';
import { styled } from '@mui/material/styles';
import ReportProblemIcon from '@mui/icons-material/ReportProblem';
import High from '@mui/icons-material/PriorityHigh';
import DoneIcon from '@mui/icons-material/Done';

import {
  GridRenderCellParams,
} from '@mui/x-data-grid';

export const STATUS_OPTIONS = ['Medio', 'OK', 'Crítico'];

interface StatusProps {
  status: string;
}

const StyledChip = styled(Chip)(({ theme }) => ({
  justifyContent: 'left',
  '& .icon': {
    color: 'inherit',
  },
  '&.OK': {
    color: (theme.vars || theme).palette.success.dark,
    border: `1px solid ${(theme.vars || theme).palette.success.main}`,
  },
  '&.Medio': {
    color: (theme.vars || theme).palette.warning.dark,
    border: `1px solid ${(theme.vars || theme).palette.warning.main}`,
  },
  '&.Crítico': {
    color: (theme.vars || theme).palette.error.dark,
    border: `1px solid ${(theme.vars || theme).palette.error.main}`,
  },
}));

const Status = React.memo((props: StatusProps) => {
  const { status } = props;

  let icon: any = null;
  if (status === 'Crítico') {
    icon = <ReportProblemIcon className="icon" />;
  } else if (status === 'Medio') {
    icon = <High className="icon" />;
  } else if (status === 'OK') {
    icon = <DoneIcon className="icon" />;
  }

  let label: string = status;
  if (status === 'Medio') {
    label = 'Medio';
  }

  return (
    <StyledChip
      className={status}
      icon={icon}
      size="small"
      label={label}
      variant="outlined"
      onClick={() => {
        /* Funcion vacia por ahora para evitar issue de que boton no tiene onClick*/
      }}
    />
  );
});

export function renderStatus(params: GridRenderCellParams<any, string>) {
  if (params.value == null) {
    return '';
  }

  return <Status status={params.value} />;
}