<script lang="ts">
  import { Router, Route } from "svelte-navigator";
  import Calculator from "./components/Calculator.svelte";
</script>

<Router>
  <main>
    <div class="retro-background"></div>
    <div class="calculator-container">
      <Route path="/">
        <Calculator />
      </Route>
    </div>
  </main>
</Router>

<style>
  main {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: #f0e4d4;
    position: relative;
    overflow: hidden;
  }
  
  .retro-background {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: 
      radial-gradient(circle at 20% 30%, rgba(255, 210, 170, 0.4) 0%, transparent 30%),
      radial-gradient(circle at 80% 70%, rgba(180, 220, 210, 0.4) 0%, transparent 30%),
      repeating-linear-gradient(45deg, rgba(0, 0, 0, 0.03) 0px, rgba(0, 0, 0, 0.03) 2px, transparent 2px, transparent 4px);
    z-index: 0;
    opacity: 0.7;
  }
  
  .calculator-container {
    position: relative;
    z-index: 1;
    filter: drop-shadow(0 20px 50px rgba(0, 0, 0, 0.3));
  }
</style> 