<script>
  // State variable to track if the menu is open or closed
  let isOpen = true;

  // Function to toggle the menu state
  // function toggleMenu() {
  //   isOpen = !isOpen;
  // }
</script>

<!-- You can toggle the class to show or hide the menu -->
<div class="side-menu" class:open={isOpen}>
  <!-- Side menu content here -->
  <!-- <button on:click={toggleMenu}>Close Menu</button> -->
  <ul>
    <li>Home</li>
    <li>About</li>
    <li>Contact</li>
  </ul>
</div>
<!-- 
<button on:click={toggleMenu}>Open Menu</button> -->

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
</style>