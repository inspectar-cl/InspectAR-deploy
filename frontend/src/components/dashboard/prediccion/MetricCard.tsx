"use client";
import React from "react";
import { Card, Avatar, Typography, Box, LinearProgress } from "@mui/material";

interface MetricCardProps {
  title: string;
  value: number;
  icon: React.ReactNode;
  color: string;
  suffix?: string;
}

function MetricCard({ title, value, icon, color, suffix = "" }: MetricCardProps): React.ReactElement {
  return <Card
    sx={{
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      justifyContent: "space-between",
      height: 230,
      width: 210,
      p: 2,
      borderRadius: 3,
      boxShadow: 3,
      //backgroundColor: `${color}22`,
    }}
  >
    <Avatar sx={{ bgcolor: color, height: 48, width: 48 }}>{icon}</Avatar>
    <Typography variant="subtitle2" sx={{ mt: 1, color: "text.secondary" }}>
      {title}
    </Typography>
    <Typography variant="h4" sx={{ mb: 2 }}>
      {value.toFixed(0)}{suffix}
    </Typography>
    <Box sx={{ flexGrow: 1 }} />
    <Box sx={{ display: "flex", justifyContent: "center", mb: 2 }}>
      <LinearProgress
        variant="determinate"
        value={Math.min(value, 100)}
        sx={{
          height: 10,
          width: 120,
          borderRadius: 5,
          backgroundColor: "#e0e0e0",
          "& .MuiLinearProgress-bar": { backgroundColor: color },
        }}
      />
    </Box>
  </Card>
}

export default MetricCard;