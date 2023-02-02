<script lang="ts">
  import { onMount } from "svelte";
  import LazyDelay from "./LazyDelay.svelte";

  export let component: () => Promise<any>;
  export let delayMs = null;

  let loadedComponent = null;

  let props: { [index: string]: any };
  $: {
    const { component, delayMs, ...restProps } = $$props;
    props = restProps;
  }

  onMount(() => {
    component().then((module) => {
      loadedComponent = module.default;
    });
  });
</script>

{#if loadedComponent}
  <svelte:component this={loadedComponent} {...props} />
{:else}
  <LazyDelay {delayMs}>
    <slot />
  </LazyDelay>
{/if}
