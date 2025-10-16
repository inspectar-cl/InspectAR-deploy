export const severidadColor = (sev: string): string =>
  sev === "Alta" ? "error" : sev === "Media" ? "warning" : "success";

export const getColorBySeverity = (sev: string): string => {
  if (sev === "Alta") return "#ef5350"; // rojo
  if (sev === "Media") return "#ffb300"; // amarillo
  return "#66bb6a"; // verde
};
