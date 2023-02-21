<script lang="ts">
  import {onDestroy, onMount} from "svelte/internal";
  import Check from "~/icons/Check.svelte";
  import Cross from "~/icons/Cross.svelte";
  import {fetchStatusOfService, serviceGraphUrl} from "~/utils/api";
  import type {BeanStatus, Service} from "~/utils/structures";
  import LazyDelay from "./LazyDelay.svelte";
  import ResponseTimeGraph from "./ResponseTimeGraph.svelte";

  let timeInterval;

  onMount(() => {
    renderBeans(days);
    timeInterval = setInterval(() => renderBeans(days), 10000);
  });

  onDestroy(() => {
    clearInterval(timeInterval);
  });

  export let service: Service;
  export let days: number;

  let graphWidth: number = 100;
  let beanStatus: BeanStatus;

  async function renderBeans(d: number) {
    beanStatus = (await fetchStatusOfService(service.id, days)) as BeanStatus;
  }
</script>

<div class="service-item">
  <div class="title-row">
    <h2>{service.name}</h2>
    {#if beanStatus}
      {#if beanStatus.current.state == 1}
        <div class="card-badge check-icon">
          <div>
            <Check size={14} />
          </div>
          <div>ONLINE</div>
        </div>
      {:else if beanStatus.current.state == 2}
        <div class="card-badge cross-icon">
          <div>
            <Cross size={14} />
          </div>
          <div>OFFLINE</div>
        </div>
      {/if}
    {/if}
  </div>
  <div class="beans">
    {#if beanStatus}
      {#each Array(90 - beanStatus.beans.length) as _}
        <div class="bean-view bean-unknown" />
      {/each}
      {#each beanStatus.beans as bean}
        <div class="bean-view {bean.state == 1 ? 'bean-on' : 'bean-off'}" />
      {/each}
    {:else}
      <LazyDelay delayMs={500} />
    {/if}
  </div>
  <div class="bean-timeline">
    <div>{days} days ago</div>
    <div>Today</div>
  </div>
  <div class="graph" bind:clientWidth={graphWidth}>
    <ResponseTimeGraph width={graphWidth} dataUrl={serviceGraphUrl(service.id)} />
  </div>
</div>

<style lang="scss">
  .service-item {
    display: flex;
    flex-direction: column;
    margin-top: 1.5rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #21242d;

    &:last-child {
      border-bottom: none;
    }

    .title-row {
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: space-between;

      h2 {
        margin: 0;
        line-height: 150%;
        font-size: 20px;
      }
    }

    .beans {
      display: grid;
      margin-top: 0.75rem;
      grid-template-columns: repeat(90, 1fr);
      gap: 1px;

      @media screen and (max-width: 400px) {
        & {
          grid-template-columns: repeat(30, 1fr);

          .bean-view:nth-last-child(30) {
            border-radius: 4px 1px 1px 4px;
          }

          .bean-view:nth-last-child(n + 31) {
            display: none;
          }
        }
      }

      @media screen and (min-width: 400px) and (max-width: 600px) {
        & {
          grid-template-columns: repeat(60, 1fr);

          .bean-view:nth-last-child(60) {
            border-radius: 4px 1px 1px 4px;
          }

          .bean-view:nth-last-child(n + 61) {
            display: none;
          }
        }
      }

      @media screen and (min-width: 600px) {
        & {
          grid-template-columns: repeat(90, 1fr);

          .bean-view:nth-last-child(90) {
            border-radius: 4px 1px 1px 4px;
          }

          .bean-view:nth-last-child(n + 91) {
            display: none;
          }
        }
      }

      .bean-view {
        height: 32px;
        border-radius: 1px;

        &:first-child {
          border-radius: 4px 1px 1px 4px;
        }

        &:last-child {
          border-radius: 1px 4px 4px 1px;
        }

        &.bean-unknown {
          background-color: var(--bean-unknown);

          &:hover {
            background-color: var(--bean-unknown-hover);
          }
        }

        &.bean-on {
          background-color: var(--bean-green);

          &:hover {
            background-color: var(--bean-green-hover);
          }
        }

        &.bean-off {
          background-color: var(--bean-red);

          &:hover {
            background-color: var(--bean-red-hover);
          }
        }
      }
    }

    .bean-timeline {
      display: flex;
      justify-content: space-between;
      margin-top: 0.75rem;
      margin-inline: 0.25rem;
      font-size: 13px;
      color: #8a91a5;
    }

    .card-badge {
      display: flex;
      border-radius: 8px;
      padding: 5px 7px;
      font-size: 10pt;
      background: #35425c;
      text-transform: uppercase;
      gap: 4px;
      align-items: center;

      div {
        display: flex;
        align-items: center;
      }

      &.check-icon {
        background-color: var(--bean-green);
      }

      &.cross-icon {
        background-color: var(--bean-red);
      }
    }
  }
</style>
