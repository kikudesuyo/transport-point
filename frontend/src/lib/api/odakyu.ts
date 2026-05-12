import type { UnifiedPoint } from '../cache';

export type OdakyuPoint = UnifiedPoint;

export async function fetchOdakyu(): Promise<OdakyuPoint[]> {
    const res = await fetch('/api/v1/odakyu');
    if (!res.ok) throw new Error('Odakyu の取得に失敗しました');
    return res.json();
}
