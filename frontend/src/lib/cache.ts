type ExpiryInfo = { points: number; date: string };
type SubPoint = { name: string; balance: number };
export type UnifiedPoint = { 
    provider: string; 
    balance: number; 
    expiry_date: string; 
    expiry_list: ExpiryInfo[]; 
    hasError?: boolean; 
    sub_points?: SubPoint[] 
};

export type CacheData = {
    data: UnifiedPoint[];
    total: number;
    timestamp: number;
};

const CACHE_KEY = 'point_hub_cache';

export function saveCache(data: UnifiedPoint[], total: number) {
    const cacheData: CacheData = {
        data,
        total,
        timestamp: Date.now()
    };
    localStorage.setItem(CACHE_KEY, JSON.stringify(cacheData));
    return cacheData.timestamp;
}

export function loadCache(): CacheData | null {
    if (typeof localStorage === 'undefined') return null;
    const cached = localStorage.getItem(CACHE_KEY);
    if (!cached) return null;
    try {
        return JSON.parse(cached) as CacheData;
    } catch (e) {
        console.error('Failed to parse cache:', e);
        return null;
    }
}

export function clearCache() {
    localStorage.removeItem(CACHE_KEY);
}
