import type { UnifiedPoint } from '../cache';

export type ToeiPoint = UnifiedPoint;

export async function fetchToei(): Promise<ToeiPoint[]> {
    const res = await fetch((import.meta.env.VITE_API_BASE || '') + '/api/v1/toei');
    if (!res.ok) throw new Error('Toei の取得に失敗しました');
    return res.json();
}
