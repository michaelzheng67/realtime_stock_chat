<script>
	import Modal from './Modal.svelte';
	import { onMount, onDestroy } from 'svelte';

	let showModal = false;
	let chat = [];
	let websocket;

	function connect() {
		websocket = new WebSocket("ws://localhost:9000/trader-ws");

		websocket.onopen = function(event) {
		console.log("Connected to WebSocket");
		};

		websocket.onmessage = function(event) {
		// Push the new message to the messages array
		chat = [...chat, event.data];
		};

		websocket.onerror = function(error) {
		console.error("WebSocket Error:", error);
		};

		websocket.onclose = function(event) {
		console.log("WebSocket connection closed", event);
		};
	}

	onMount(() => {
		connect();
	});

	onDestroy(() => {
		if (websocket) {
		websocket.close();
		}
	});
</script>

<button on:click={() => (showModal = true)}> Chat </button>

<Modal bind:showModal>
	<h2 slot="header">
		Live Chat
	</h2>

	<div>
	{#each chat as string}
		<p>{string}</p>
	{/each}
	</div>

	<input type="text">
</Modal>