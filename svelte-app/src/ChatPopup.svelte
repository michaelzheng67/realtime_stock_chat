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

	let inputMessage = '';

    // Existing WebSocket and sendMessage function...

    function handleEnterPress(event) {
        if (event.key === 'Enter') {
            sendMessage(inputMessage);
            inputMessage = ''; // Clear the input field after sending the message
        }
    }

	// send websocket message
	function sendMessage(message) {
		if (websocket && websocket.readyState === WebSocket.OPEN) {
			websocket.send(message);
		} else {
			console.error("WebSocket is not connected.");
		}
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

	<div class="scrollable-content">
	{#each chat as string}
		<p>{string}</p>
	{/each}
	</div>

	<input type="text" bind:value={inputMessage} on:keydown={handleEnterPress}>
</Modal>



<style>

.scrollable-content {
    height: 300px; /* or max-height: 300px; */
    overflow: auto; /* or overflow: scroll; */
}
</style>
