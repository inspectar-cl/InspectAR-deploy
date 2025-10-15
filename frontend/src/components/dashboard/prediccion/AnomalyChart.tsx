/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
"use client";
import React from "react";
import dynamic from "next/dynamic";
import { Prediccion } from "@/types/prediccion";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

export const AnomalyChart: React.FC<{ predicciones: Prediccion[]; activoNombre: string }> = ({
  predicciones,
  activoNombre,
}) => {
  const series = [
    { name: "Anomaly Score", data: predicciones.map(p => ({ x: p.timestamp, y: p.anomalyScore })) },
    { name: "Anomaly Likelihood", data: predicciones.map(p => ({ x: p.timestamp, y: p.anomalyLikelihood })) },
  ];

  const options = {
    chart: { type: "line", zoom: { enabled: true } },
    xaxis: { type: "datetime" },
    yaxis: [{ title: { text: "Anomaly Score" } }, { opposite: true, title: { text: "Likelihood" } }],
    tooltip: { shared: true, intersect: false },
    title: { text: `Predicciones para ${activoNombre}`, align: "left" },
  };

  return <Chart options={options} series={series} type="line" height={300} />;
};

export default AnomalyChart;