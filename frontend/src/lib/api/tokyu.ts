import type { UnifiedPoint } from '../cache';

export type TokyuPoint = UnifiedPoint;

export async function fetchTokyu(): Promise<TokyuPoint[]> {
    const res = await fetch('/api/v1/tokyu');
    if (!res.ok) throw new Error('Tokyu の取得に失敗しました');
    return res.json();
}
