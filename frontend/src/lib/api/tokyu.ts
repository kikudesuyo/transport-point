import type { UnifiedPoint } from '../cache';

export type TokyuPoint = UnifiedPoint;

export async function fetchTokyu(): Promise<TokyuPoint[]> {
    const res = await fetch((import.meta.env.VITE_API_BASE || '') + '/api/v1/tokyu');
    if (!res.ok) throw new Error('Tokyu の取得に失敗しました');
    return res.json();
}
