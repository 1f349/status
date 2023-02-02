<script lang="ts">
  import { onMount } from "svelte";

  export let delayMs = null;

  let timeout: string | number | NodeJS.Timeout;
  let showFallback = !delayMs;

  onMount(() => {
    if (delayMs) {
      timeout = setTimeout(() => {
        showFallback = true;
      }, delayMs);
    }
    return () => clearTimeout(timeout);
  });
</script>

{#if showFallback}
  <slot />
{/if}
