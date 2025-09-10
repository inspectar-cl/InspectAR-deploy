import type { ReactElement } from "react";

export interface Product {
  id: string;
  icon?: ReactElement;
  name: string;
  updatedAt: Date;
}