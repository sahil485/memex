// Wails Runtime type definitions
export function EventsOn(eventName: string, callback: (...data: any) => void): void;
export function EventsOff(eventName: string): void;
export function EventsEmit(eventName: string, ...data: any): void;
export function LogPrint(message: string): void;
export function LogTrace(message: string): void;
export function LogDebug(message: string): void;
export function LogInfo(message: string): void;
export function LogWarning(message: string): void;
export function LogError(message: string): void;
export function LogFatal(message: string): void;
export function Quit(): void;
export function WindowHide(): void;
