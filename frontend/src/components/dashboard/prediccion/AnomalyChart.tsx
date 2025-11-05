"use client";
import React from "react";
import dynamic from "next/dynamic";
import { type Prediccion } from "@/types/prediccion";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

function AnomalyChart({ predicciones, activoNombre }: { predicciones: Prediccion[]; activoNombre: string }): React.ReactElement {
  const series = [
    { name: "Anomaly Score", data: predicciones.map(p => ({ x: p.timestamp, y: p.anomalyScore })) },
    { name: "Anomaly Likelihood", data: predicciones.map(p => ({ x: p.timestamp, y: p.anomalyLikelihood })) },
    { name: "Threshold", data: predicciones.map(p => ({ x: p.timestamp, y: p.threshold })) },
  ];

  const options = {
    chart: { type: "line" as const, zoom: { enabled: true } },
    xaxis: { type: "datetime" as const },
    yaxis: [{ title: { text: "Anomaly Score" } }, { opposite: true, title: { text: "Likelihood" } }],
    tooltip: { 
      shared: true, 
      intersect: false,
      theme: 'dark'
    },
    title: { text: `Predicciones para ${activoNombre}`, align: "left" as const },
  };

  return <Chart options={options} series={series} type="line" height={300} />;
};

export default AnomalyChart;