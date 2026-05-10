<script lang="ts">
	let { errors } = $props<{ errors: string[] }>();
</script>

{#if errors.length > 0}
	<div class="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700 shadow-sm max-w-3xl mx-auto mb-10">
		<div class="flex items-center gap-2 mb-2">
			<svg class="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
			<h2 class="text-base font-bold">一部の取得に失敗しました</h2>
		</div>
		<ul class="list-disc ml-7 text-sm text-red-600">
			{#each errors as error}
				<li>{error}</li>
			{/each}
		</ul>
	</div>
{/if}
