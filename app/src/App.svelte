<script lang="ts">
	import Loading from "./Loading/Component.svelte";
	import Alert from "./Alert/Component.svelte";
	import MyApps from "./UI/MyApps.svelte";
	import { route } from "./Router";
	import Status from "./Status";
	import OnLoaded from "./OnLoaded.svelte";
	import OauthApp from "./UI/OauthApp.svelte";

	let loaded = false;
	Status.Login().then(() => {
		loaded = true;
	});
</script>

<main>
	{#if loaded}
		<OnLoaded />
		{#if $route.INFO}
			<MyApps />
		{/if}
		{#if $route.APP}
			<OauthApp />
		{/if}
	{/if}
	<Alert />
	<Loading />
</main>

<style>
	main {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		text-align: center;
		padding: 1em;
		width: 500px;
		margin: auto;
	}
</style>
