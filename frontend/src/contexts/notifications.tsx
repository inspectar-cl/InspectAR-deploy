/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
'use client';

import * as React from 'react';

export interface Notification {
    id: string; // uuid
    title: string; // ej: "Bomba A — Sensor de Presión inactivo"
    body?: string; // detalle opcional
    severity?: 'info'|'warning'|'error'|'success';
    ts: Date;
    read?: boolean;
    meta?: Record<string, unknown>;// sensorId, assetName, etc.
}


export interface NotificationsContextType {
    notifications: Notification[];
    unreadCount: number;
    push: (n: Omit<Notification, 'id'|'ts'|'read'> & { id?: string; ts?: Date; read?: boolean }) => void;
    markAllRead: () => void;
    remove: (id: string) => void;
    clear: () => void;
// Registro de estado conocido por sensor para detectar transiciones
    getSensorStatus: (sensorId: string) => 'connected'|'disconnected'|'never_connected'|undefined;
    setSensorStatus: (sensorId: string, status: 'connected'|'disconnected'|'never_connected') => void;
}


const NotificationsContext = React.createContext<NotificationsContextType | undefined>(undefined);


function uid() {
    return crypto?.randomUUID?.() ?? Math.random().toString(36).slice(2);
}


export function NotificationsProvider({ children }: { children: React.ReactNode }) {
    const [notifications, setNotifications] = React.useState<Notification[]>([]);
    const sensorStatusRef = React.useRef<Map<string, 'connected'|'disconnected'|'never_connected'>>(new Map());
    const push: NotificationsContextType['push'] = React.useCallback((n) => {
        setNotifications((prev) => [
            { id: n.id ?? uid(), ts: n.ts ?? new Date(), read: n.read ?? false, severity: 'warning', ...n },
            ...prev,
        ]);
    }, []);


    const markAllRead = React.useCallback(() => {
        setNotifications((prev) => prev.map((x) => ({ ...x, read: true })));
    }, []);


    const remove = React.useCallback((id: string) => {
        setNotifications((prev) => prev.filter((x) => x.id !== id));
    }, []);


    const clear = React.useCallback(() => {
        setNotifications([]);
    }, []);


    const getSensorStatus = React.useCallback((sensorId: string) => {
        return sensorStatusRef.current.get(sensorId);
    }, []);


    const setSensorStatus = React.useCallback((sensorId: string, status: 'connected'|'disconnected'|'never_connected') => {
        sensorStatusRef.current.set(sensorId, status);
    }, []);


    const unreadCount = notifications.reduce((acc, n) => acc + (n.read ? 0 : 1), 0);


    const value: NotificationsContextType = React.useMemo(() => ({
        notifications,
        unreadCount,
        push,
        markAllRead,
        remove,
        clear,
        getSensorStatus,
        setSensorStatus,
    }), [notifications, unreadCount, push, markAllRead, remove, clear, getSensorStatus, setSensorStatus]);


    return (
        <NotificationsContext.Provider value={value}>{children}</NotificationsContext.Provider>
    );
}


export function useNotifications() {
    const ctx = React.useContext(NotificationsContext);
    if (!ctx) throw new Error('useNotifications must be used within NotificationsProvider');
        return ctx;
}