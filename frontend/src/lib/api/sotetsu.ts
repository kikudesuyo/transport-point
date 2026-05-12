import type { UnifiedPoint } from '../cache';

export type SotetsuPoint = UnifiedPoint;

export async function fetchSotetsu(): Promise<SotetsuPoint[]> {
    const res = await fetch('/api/v1/sotetsu');
    if (!res.ok) throw new Error('Sotetsu の取得に失敗しました');
    return res.json();
}
