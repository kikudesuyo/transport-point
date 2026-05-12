import type { UnifiedPoint } from '../cache';

export type KeikyuPoint = UnifiedPoint;

export async function fetchKeikyu(): Promise<KeikyuPoint[]> {
    const res = await fetch('/api/v1/keikyu');
    if (!res.ok) throw new Error('Keikyu の取得に失敗しました');
    return res.json();
}
