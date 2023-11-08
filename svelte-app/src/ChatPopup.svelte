<script>
	import Modal from './Modal.svelte';

	let showModal = false;
</script>

<button on:click={() => (showModal = true)}> show modal </button>

<Modal bind:showModal>
	<h2 slot="header">
		Live Chat
	</h2>

	<ol class="definition-list">
		<li>of or relating to modality in logic</li>
		<li>
			containing provisions as to the mode of procedure or the manner of taking effect —used of a
			contract or legacy
		</li>
		<li>of or relating to a musical mode</li>
		<li>of or relating to structure as opposed to substance</li>
		<li>
			of, relating to, or constituting a grammatical form or category characteristically indicating
			predication
		</li>
		<li>of or relating to a statistical mode</li>
	</ol>

	<input type="text">
</Modal>