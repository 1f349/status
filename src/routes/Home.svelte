<script lang="ts">
  import {onDestroy, onMount} from "svelte/internal";
  import LazyDelay from "~/lib/LazyDelay.svelte";
  import ServiceView from "~/lib/ServiceView.svelte";
  import {fetchStatus} from "~/utils/api";
  import type {Group} from "~/utils/structures";

  let timeInterval;

  onMount(() => {
    renderBeans();
    timeInterval = setInterval(() => renderBeans(), 10000);
  });

  onDestroy(() => {
    clearInterval(timeInterval);
  });

  async function renderBeans() {
    groups = (await fetchStatus()) as Group[];
  }

  let innerWidth: number = window.innerWidth;
  $: days = calculateNumberOfDays(innerWidth);

  let groups: Group[] = [];

  function calculateNumberOfDays(w: number) {
    if (w >= 600) return 90;
    if (w >= 400) return 60;
    return 30;
  }
</script>

<svelte:window bind:innerWidth />

<div class="wrapper-block">
  {#if groups}
    {#each groups as y (y.id)}
      <div class="group-item rounded-card">
        <h1>{y.name}</h1>
        {#each y.services as z (z.id)}
          <ServiceView service={z} {days} />
        {/each}
      </div>
    {/each}
  {:else}
    <LazyDelay delayMs={500} />
  {/if}
</div>

<style lang="scss">
  .wrapper-block {
    display: flex;
    gap: 12px;
    flex-direction: column;
    margin: 12px;

    .group-item {
      padding: 0.75rem;
      border-radius: 12px;
      box-shadow: 0 4px 7px 0 #00000008;
      border: 1px solid #21242d;

      h1 {
        color: #9c9487;
        margin: 10px 0 0 10px;
        font-size: 16px;
        font-weight: 500;
        padding-left: 0;
        text-transform: uppercase;
      }
    }
  }
</style>
