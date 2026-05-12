<script lang="ts">
	import { onMount } from 'svelte';
	import Hero from '$lib/components/Hero.svelte';
	import ProviderCard from '$lib/components/ProviderCard.svelte';
	import ErrorAlert from '$lib/components/ErrorAlert.svelte';
	import { loadCache, saveCache, type UnifiedPoint } from '$lib/cache';

	let details: UnifiedPoint[] = $state([]);
	let totalBalance: number = $state(0);
	let isLoading = $state(false);
	let errors: string[] = $state([]);
	let lastUpdated: number | null = $state(null);

	const providers = ['tokyu', 'metpo', 'toei', 'sotetsu', 'keikyu', 'odakyu', 'tobu'];

	onMount(() => {
		const cached = loadCache();
		if (cached) {
			details = cached.data;
			totalBalance = cached.total;
			lastUpdated = cached.timestamp;
		}
	});

	async function fetchPoints() {
		isLoading = true;
		errors = [];
		details = [];
		totalBalance = 0;

		const fetchPromises = providers.map(async (provider) => {
			try {
				const res = await fetch(`/api/v1/${provider}`);
				if (!res.ok) {
					throw new Error(`${provider} の取得に失敗しました`);
				}
				const data: UnifiedPoint[] = await res.json();
				
				if (data && data.length > 0) {
					details = [...details, ...data];
					totalBalance += data.reduce((acc, curr) => acc + curr.balance, 0);
				}
			} catch (e: any) {
				errors = [...errors, e.message || `${provider} で不明なエラーが発生しました`];
				const providerName = provider.charAt(0).toUpperCase() + provider.slice(1);
				details = [...details, { provider: providerName, balance: 0, expiry_date: '--', expiry_list: [], hasError: true }];
			}
		});

		await Promise.allSettled(fetchPromises);
		
		lastUpdated = saveCache(details, totalBalance);
		isLoading = false;
	}
</script>

<svelte:head>
	<title>Point Hub Dashboard</title>
</svelte:head>

<main class="min-h-screen bg-gray-50 text-gray-900 font-sans">
	<div class="max-w-5xl mx-auto px-6 py-16">
		<Hero onSync={fetchPoints} {isLoading} {lastUpdated} />
		<ErrorAlert {errors} />

		{#if details.length > 0 || isLoading}
			<div class="animate-in fade-in slide-in-from-bottom-4 duration-500">
				<!-- Total Balance Card -->
				<div class="mb-10 p-8 rounded-2xl bg-white border border-gray-200 shadow-sm flex flex-col items-center">
					<h2 class="text-sm text-gray-500 font-medium mb-2 tracking-wide uppercase">Total Balance</h2>
					<div class="text-6xl md:text-7xl font-bold text-gray-900 tracking-tight">
						{totalBalance.toLocaleString()} <span class="text-3xl text-gray-400 font-normal ml-1">pt</span>
					</div>
				</div>

				<!-- Details Grid -->
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each details as detail (detail.provider)}
						<ProviderCard {detail} />
					{/each}
				</div>
			</div>
		{/if}
	</div>
</main>
