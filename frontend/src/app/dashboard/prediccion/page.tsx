/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
import type { Metadata } from 'next';
import { config } from '@/config';
import PrediccionClient from './calls';

export const metadata = { title: `Predicciones | Dashboard | ${config.site.name}` } satisfies Metadata;

export default function PrediccionPage() {
  return <PrediccionClient />;
}