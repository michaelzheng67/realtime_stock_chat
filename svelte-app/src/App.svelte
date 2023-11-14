<script>
	import SideMenu from './SideMenu.svelte';
	import ChatPopup from './ChatPopup.svelte';
	import { onMount, onDestroy } from 'svelte';

	let websocket;
	let page = "home";
	let stocks = {
		AAPL: '119',
		MSFT: '213',
		AMZN: '3116',
		GOOG: '1735',
		META: '276',
		TSLA: '408',
		BABA: '309',
	}

	let colors = {
		AAPL: 'grey',
		MSFT: 'grey',
		AMZN: 'grey',
		GOOG: 'grey',
		META: 'grey',
		TSLA: 'grey',
		BABA: 'grey',
	}

	function updatePage(newPage) {
		page = newPage
	
	}

	function connect() {
		websocket = new WebSocket("ws://localhost:8000/stock-ws");

		websocket.onopen = function(event) {
		console.log("Connected to WebSocket");
		};

		websocket.onmessage = function(event) {

		let parsed
		let sym 
		let c 

		parsed = JSON.parse(event.data)

		sym = parsed.sym
		c = parsed.c


		if (sym in stocks) {
			// update stock movement color
			if (stocks[sym] > c) {
				colors[sym] = 'red'
			} else if (stocks[sym] < c) {
				colors[sym] = 'green'
			} else {
				colors[sym] = 'grey'
			}

			// update price
			stocks[sym] = parseFloat(c).toFixed(2)
		}
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

<main>

	<div class="header">
		{#if page == "home"}
		<h1>Real Time Stock Chat</h1>
		<h3>View stock data in real time and chat to your fellow traders</h3>

		{:else}
		<h1>{page}</h1>
		<h3>View stock data in real time and chat to your fellow traders</h3>
		{/if}

		
	</div>

	<SideMenu {stocks} {updatePage} {colors}/>
	<div class="bottom-right">
		<ChatPopup />
	</div>

	
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
		display: flex;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
		margin-bottom: 5px;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}

	.header {
		text-align: left;
		flex: 1;
		margin-left: 300px;
	}

	.bottom-right {
		position: fixed; /* Fixed positioning relative to the viewport */
		right: 0; /* Align to the right side of the viewport */
		bottom: 0; /* Align to the bottom of the viewport */
		margin-right: 10px; /* Add some space from the right edge of the viewport */
		margin-bottom: 10px; /* Add some space from the bottom edge of the viewport */
		/* Additional styling */
		padding: 10px;
		/* background-color: #f8f9fa; */
		border-radius: 5px;
		/* box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2); */
		z-index: 1000; /* Make sure it's above other elements */
  	}
</style>