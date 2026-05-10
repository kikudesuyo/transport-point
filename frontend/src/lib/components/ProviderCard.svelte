<script lang="ts">
	type ExpiryInfo = { points: number; date: string };
	type UnifiedPoint = { provider: string; balance: number; expiry_date: string; expiry_list: ExpiryInfo[] };
	
	let { detail } = $props<{ detail: UnifiedPoint }>();
</script>

<div class="rounded-xl bg-white border border-gray-200 p-6 shadow-sm flex flex-col justify-between hover:shadow-md transition-shadow">
	<div>
		<h3 class="text-lg font-semibold text-gray-700 mb-4">{detail.provider}</h3>
		<div class="text-3xl font-bold text-gray-900 mb-2">
			{detail.balance.toLocaleString()} <span class="text-lg text-gray-500 font-medium">pt</span>
		</div>
	</div>
	
	{#if detail.expiry_date && detail.expiry_date !== "--"}
		<div class="mt-6 pt-4 border-t border-gray-100 text-sm">
			<div class="flex items-center justify-between mb-2">
				<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-50 text-red-600 border border-red-100">
					最短有効期限
				</span>
				<span class="font-medium text-gray-800">{detail.expiry_date}</span>
			</div>
			
			{#if detail.expiry_list && detail.expiry_list.length > 0}
				<div class="space-y-1.5 mt-3">
					{#each detail.expiry_list as exp}
						<div class="flex justify-between text-xs items-center p-1.5 rounded bg-gray-50">
							<span class="text-gray-500">{exp.date}</span>
							<span class="font-semibold text-gray-700">{exp.points.toLocaleString()} pt</span>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{:else}
		<div class="mt-6 pt-4 border-t border-gray-100 text-sm text-gray-400 italic">
			有効期限のあるポイントはありません
		</div>
	{/if}
</div>
