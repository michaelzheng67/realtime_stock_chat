<script>
  // State variable to track if the menu is open or closed
  let isOpen = true;
  export let stocks;

</script>
<!-- You can toggle the class to show or hide the menu -->
<div class="side-menu" class:open={isOpen}>

  <!-- <ul>
    <li>Home</li>
    <li>About</li>
    <li>Contact</li>
  </ul> -->

  <!-- <ul> -->
    {#each Object.entries(stocks) as [key, value]}
      <!-- <li>{key}: {value}</li> -->
      <div>
        <h3>{key} : ${value}</h3>
      </div>

    {/each}
  <!-- </ul> -->
</div>


<style>
  .side-menu {
    width: 250px;
    position: fixed;
    left: -250px;
    top: 0;
    bottom: 0;
    background: #333;
    color: white;
    transition: left 0.3s;
    padding: 1em;
  }

  .side-menu.open {
    left: 0;
  }

  li {
    margin-bottom: 10px; /* Adjust the value as needed */
  }
</style>