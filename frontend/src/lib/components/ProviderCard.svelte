<script lang="ts">
	import type { UnifiedPoint } from '$lib/cache';

	let { detail } = $props<{ detail: UnifiedPoint }>();
</script>

<div
	class="flex flex-col justify-between rounded-xl border border-gray-200 bg-white p-6 shadow-sm transition-shadow hover:shadow-md"
>
	<div>
		<h3 class="mb-4 text-lg font-semibold text-gray-700">{detail.provider}</h3>

		{#if detail.hasError}
			<div
				class="mt-4 mb-2 rounded border border-red-100 bg-red-50 px-3 py-2 text-sm font-bold text-red-500"
			>
				取得できませんでした
			</div>
		{:else}
			<div class="mb-2 text-3xl font-bold text-gray-900">
				{detail.balance.toLocaleString()} <span class="text-lg font-medium text-gray-500">pt</span>
			</div>
			{#if detail.sub_points && detail.sub_points.length > 0}
				<div class="mt-2 mb-4 flex flex-wrap gap-x-4 gap-y-1 text-sm text-gray-600">
					{#each detail.sub_points as sp}
						<div class="flex items-center gap-1.5 rounded-full bg-gray-50 px-2.5 py-1 border border-gray-100">
							<span class="font-medium">{sp.name}:</span>
							<span>{sp.balance.toLocaleString()} pt</span>
						</div>
					{/each}
				</div>
			{/if}
		{/if}
	</div>

	{#if !detail.hasError}
		{#if detail.expiry_date && detail.expiry_date !== '--'}
			<div class="mt-6 border-t border-gray-100 pt-4 text-sm">
				<div class="mb-2 flex items-center justify-between">
					<span
						class="inline-flex items-center rounded border border-red-100 bg-red-50 px-2 py-0.5 text-xs font-medium text-red-600"
					>
						最短有効期限
					</span>
					<span class="font-medium text-gray-800">{detail.expiry_date}</span>
				</div>

				{#if detail.expiry_list && detail.expiry_list.length > 0}
					<div class="mt-3 space-y-1.5">
						{#each detail.expiry_list as exp}
							<div class="flex items-center justify-between rounded bg-gray-50 p-1.5 text-xs">
								<span class="text-gray-500">{exp.date}</span>
								<span class="font-semibold text-gray-700">{exp.points.toLocaleString()} pt</span>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{:else}
			<div class="mt-6 border-t border-gray-100 pt-4 text-sm text-gray-400 italic">
				有効期限のあるポイントはありません
			</div>
		{/if}
	{/if}
</div>
