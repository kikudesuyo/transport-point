import type { UnifiedPoint } from '../cache';

export type MetpoPoint = UnifiedPoint;

export async function fetchMetpo(): Promise<MetpoPoint[]> {
    const res = await fetch('/api/v1/metpo');
    if (!res.ok) throw new Error('Metpo の取得に失敗しました');
    return res.json();
}
