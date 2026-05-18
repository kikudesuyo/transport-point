import type { UnifiedPoint } from '../cache';

export type TobuPoint = UnifiedPoint;

export async function fetchTobu(): Promise<TobuPoint[]> {
    const res = await fetch((import.meta.env.VITE_API_BASE || '') + '/api/v1/tobu');
    if (!res.ok) throw new Error('Tobu の取得に失敗しました');
    return res.json();
}
