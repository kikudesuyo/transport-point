<script lang="ts">
	let { onSync, isLoading } = $props<{ onSync: () => void; isLoading: boolean }>();
</script>

<header class="mb-12 text-center flex flex-col items-center">
	<h1 class="text-4xl md:text-5xl font-extrabold tracking-tight mb-4 text-gray-900">
		Point Hub
	</h1>
	<p class="text-gray-500 text-base md:text-lg mb-8 max-w-xl leading-relaxed">
		交通系から日常のポイントまで。<br />散らばったポイント残高と有効期限を、ワンクリックで集計します。
	</p>
	
	<button 
		onclick={onSync} 
		disabled={isLoading}
		class="inline-flex items-center justify-center px-6 py-3 font-medium text-white transition-colors bg-blue-600 rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed shadow-sm"
	>
		{#if isLoading}
			<svg class="animate-spin -ml-1 mr-2 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			集計中...
		{:else}
			<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
			</svg>
			Sync Points
		{/if}
	</button>
</header>
